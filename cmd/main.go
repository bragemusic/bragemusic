package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/client"
	"github.com/bragemusic/core/pkg/database"
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

	aid, err := acoustid.New(scfg.AcoustIDApiKey, slogHandler)
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

	// m := mediamanager.New(slogHandler, db, "/home/lucas/dev/p/brage-player/data/music")

	// sc := server.Config{ImagePath: "/home/lucas/dev/brage-player/data/img"}
	// s := server.New(slogHandler, &m, sc)

	// logger.Info(fmt.Sprintf("serving on port %s", scfg.Port))
	// go func() {
	// 	if err = http.ListenAndServe(":"+scfg.Port, s.Handler()); err != nil {
	// 		logger.Error(err.Error())
	// 		return
	// 	}
	// }()

	syC, err := client.NewSyncer(client.Config{
		ServerBaseURL: "http://localhost:3000",
		MusicDirPath:  "/home/lucas/dev/brage/client_data/music",
		ConfigPath:    "/home/lucas/dev/brage/client_data",
		ImagePath:     "/home/lucas/dev/brage/client_data/img",
	}, slogHandler)
	if err != nil {
		panic(err)
	}
	defer syC.Close()

	err = syC.Sync(context.Background())
	if err != nil {
		panic(err)
	}
}
