package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/wiki"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
)

func main() {
	scfg, err := config.GetServerConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	scfg.DBPath = "/home/lucas/dev/brage/data/config/data.db"

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

	imp := importer.New("/home/lucas/dev/brage/importDir", "/home/lucas/dev/brage/data/music", &db, musicbrainz.MusicBrainz{}, aid, w, slogHandler)

	imp.Run(context.Background())
}
