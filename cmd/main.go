package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/trackmgr"
	"github.com/bragemusic/core/pkg/wiki"

	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	scfg, err := config.GetServerConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	dbSqlite, err := sqlx.Open("sqlite3", scfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbSqlite.Close()

	aid, err := acoustid.New(scfg.AcoustIDApiKey)
	if err != nil {
		panic(err)
	}

	db, err := database.New(dbSqlite)
	if err != nil {
		panic(err)
	}

	w := wiki.New(scfg.WikiEmail)

	tm := trackmgr.New(scfg, db, aid, w, slogHandler)

	imp := importer.New("importDir", &tm, slogHandler)
	imp.Run(context.Background())
}
