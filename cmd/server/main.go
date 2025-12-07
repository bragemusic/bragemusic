package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/server"

	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	logger := slog.New(slogHandler)

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

	m := mediamanager.New(slogHandler, db, "/home/lucas/dev/brage/data/music")

	s := server.New(slogHandler, &m, scfg)

	logger.Info(fmt.Sprintf("serving on port %d", scfg.Port))
	if err = http.ListenAndServe(fmt.Sprintf(":%d", scfg.Port), s.Handler()); err != nil {
		logger.Error(err.Error())
		return
	}
}
