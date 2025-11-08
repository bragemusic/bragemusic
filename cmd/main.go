package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/server"
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

	logger := slog.New(slogHandler)

	dbSqlite, err := sqlx.Open("sqlite3", scfg.DBPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer dbSqlite.Close()

	aid, err := acoustid.New(scfg.AcoustIDApiKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	db, err := database.New(dbSqlite)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	w := wiki.New(scfg.WikiEmail)

	tm := trackmgr.New(scfg, db, aid, w, slogHandler)
	_ = tm

	m := mediamanager.New(slogHandler, db)

	s := server.New(slogHandler, &m)

	logger.Info(fmt.Sprintf("serving on port %s", scfg.Port))
	if err = http.ListenAndServe(":"+scfg.Port, s.Handler()); err != nil {
		logger.Error(err.Error())
		return
	}

	// imp := importer.New("importDir", &tm, slogHandler)
	// imp.Run(context.Background())
}
