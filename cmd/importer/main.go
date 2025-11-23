package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/metasyncer"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/wiki"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
)

func main() {
	datadir := "/home/lucas/dev/brage/data"
	scfg, err := config.GetServerConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	scfg.DBPath = filepath.Join(datadir, "config/data.db")

	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	logger := slog.New(slogHandler)

	dbSqlite, err := sqlx.Open("sqlite3", scfg.DBPath)
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

	aid, err := acoustid.New(scfg.AcoustIDApiKey, slogHandler)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	w := wiki.New(scfg.WikiEmail)

	mb := musicbrainz.New(slogHandler)

	impCfg := importer.Config{
		ImportDirPath: filepath.Join(datadir, "..", "importDir"),
		MusicDirPath:  filepath.Join(datadir, "music"),
		ImageDirPath:  filepath.Join(datadir, "img"),
	}

	imp := importer.New(impCfg, &db, mb, aid, slogHandler)

	ms := metasyncer.New(impCfg.ImageDirPath, &db, musicbrainz.MusicBrainz{}, w, slogHandler)

	ctx := context.Background()
	imp.Run(ctx)
	ms.Sync(ctx)
}
