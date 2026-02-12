package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/filetx"
	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/metasyncer"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/wiki"

	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// FIXME: Add proper context with SIGs
	ctx := context.Background()
	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	logger := slog.New(slogHandler)

	filetx.Init(slogHandler)

	scfg, err := server.GetConfig()
	if err != nil {
		logger.Error("could not parse config", "error", err.Error())
		return
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

	w := wiki.New(scfg.Wikipedia.Email)

	mb := musicbrainz.New(slogHandler)

	imp := importer.New(impCfg, &db, mb, aid, im, slogHandler)

	ms := metasyncer.New(impCfg.ImageDirPath, &db, mb, w, im, slogHandler)

	jmgr := jobmanager.New(slogHandler, &m, &imp, &ms)

	s := server.New(slogHandler, &m, &a, &imp, &jmgr, scfg)

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
