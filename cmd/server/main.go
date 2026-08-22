package main

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bragemusic/bragemusic/assets"
	"github.com/bragemusic/bragemusic/pkg/acoustid"
	"github.com/bragemusic/bragemusic/pkg/analyser"
	"github.com/bragemusic/bragemusic/pkg/auth"
	"github.com/bragemusic/bragemusic/pkg/config"
	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/device"
	"github.com/bragemusic/bragemusic/pkg/filetx"
	"github.com/bragemusic/bragemusic/pkg/imagemagick"
	"github.com/bragemusic/bragemusic/pkg/importer"
	"github.com/bragemusic/bragemusic/pkg/internalusers"
	"github.com/bragemusic/bragemusic/pkg/jobmanager"
	"github.com/bragemusic/bragemusic/pkg/mediamanager"
	"github.com/bragemusic/bragemusic/pkg/metasyncer"
	"github.com/bragemusic/bragemusic/pkg/musicbrainz"
	"github.com/bragemusic/bragemusic/pkg/server"
	"github.com/bragemusic/bragemusic/pkg/sse"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/bragemusic/bragemusic/pkg/utils"
	"github.com/bragemusic/bragemusic/pkg/wiki"

	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	var slogHandler slog.Handler

	slogHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(slogHandler)

	filetx.Init(slogHandler)

	scfg, err := config.GetServerConfig(logger)
	if err != nil {
		logger.Error("could not parse config", "error", err.Error())
		return
	}

	logLevel, err := utils.ParseLogLevel(scfg.General.LogLevel)
	if err != nil {
		logger.Error("could not parse log level. Falling back to DEBUG", "error", err.Error())
	} else {
		switch scfg.General.LogFormat {
		case config.LogFormatPretty:
			slogHandler = tint.NewHandler(os.Stderr, &tint.Options{
				Level:      logLevel,
				TimeFormat: time.TimeOnly,
			})
		case config.LogFormatJson:
			slogHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: logLevel,
			})
		default:
			logger.Error("unknown log format", "log_format", scfg.General.LogFormat)
		}

		logger = slog.New(slogHandler)
	}

	dbPath := filepath.Join(scfg.Paths.ConfigDir, "data.db")
	dbSqlite, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer dbSqlite.Close()

	db, err := database.New(dbSqlite)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	a := auth.New(db, slogHandler)
	err = a.SetAdmin(ctx, scfg.Admin.Email, scfg.Admin.Username, scfg.Admin.Password)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	im, err := imagemagick.New(slogHandler)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	m := mediamanager.New(slogHandler, db, &im, scfg.Paths.MusicDir, scfg.Paths.ImageDir)

	impCfg := importer.Config{
		ImportDirPath:          scfg.Paths.ImportDir,
		ManualImportDirPath:    scfg.Paths.ManualImportDir,
		MusicDirPath:           scfg.Paths.MusicDir,
		ImageDirPath:           scfg.Paths.ImageDir,
		DeleteImportsOnSuccess: false,
		FinishedImportsDirPath: scfg.Paths.BackupImportDir,
	}

	im, err = imagemagick.New(slogHandler)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	aid, err := acoustid.New(scfg.AcoustID.ApiKey, slogHandler)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	distFS, err := fs.Sub(assets.DistFS, "frontend")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sseHub := sse.NewHub(&db, slogHandler)

	w := wiki.New(scfg.Wikipedia.Email, slogHandler)

	mb := musicbrainz.New(slogHandler)

	imp := importer.New(impCfg, &db, sseHub, mb, aid, im, slogHandler)

	ms := metasyncer.New(impCfg.ImageDirPath, &db, mb, w, im, internalusers.MetaSyncer, slogHandler)

	jmgr := jobmanager.New(slogHandler)
	jmgr.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobImporterRun,
		CronExpr: scfg.Jobs.Importer,
		Run:      imp.Run,
	})
	jmgr.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobMetaSyncRun,
		CronExpr: scfg.Jobs.MetaSyncer,
		Run:      ms.Sync,
	})
	jmgr.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobMediaManagerUpdateSyncItems,
		CronExpr: scfg.Jobs.SearchItems,
		Run:      m.UpdateSearchItems,
	})
	jmgr.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobAuthExpiredTokenCleanup,
		CronExpr: scfg.Jobs.TokenCleanup,
		Run:      a.TokenCleanupJob,
	})

	if scfg.Analyser.BaseURL != "" {
		ana := analyser.New(scfg.Analyser.BaseURL, scfg.Paths.MusicDir, &db, slogHandler)
		jmgr.RegisterJob(ctx, jobmanager.JobDefinition{
			Type:     types.JobAnalyserRunTrackAnalysis,
			CronExpr: scfg.Jobs.TrackAnalysis,
			Run:      ana.RunTrackAnalysisJob,
		})
	} else {
		logger.WarnContext(ctx, "no analyser base url set. No analysis will be performed")
	}

	dm := device.NewManager(sseHub, &db, slogHandler)

	s := server.New(slogHandler, &m, &a, &imp, &mb, &jmgr, sseHub, &dm, distFS, scfg)

	// logger.Info(fmt.Sprintf("serving on port %d", scfg.Port))
	// if err = http.ListenAndServe(fmt.Sprintf(":%d", scfg.Port), s.Handler()); err != nil {
	// 	logger.Error(err.Error())
	// 	return
	// }
	if err = s.Start(ctx); err != nil {
		logger.Error(err.Error())
		return
	}
}
