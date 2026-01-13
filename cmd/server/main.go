package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/server"

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

	s := server.New(slogHandler, &m, &a, scfg)

	logger.Info(fmt.Sprintf("serving on port %d", scfg.Port))
	if err = http.ListenAndServe(fmt.Sprintf(":%d", scfg.Port), s.Handler()); err != nil {
		logger.Error(err.Error())
		return
	}
}
