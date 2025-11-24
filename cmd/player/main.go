package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
)

func main() {
	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		panic(err)
	}

	ap, err := audioplayer.New(audioplayer.Config{PlayerName: "brage", MusicDirPath: "/home/lucas/dev/brage/data/music"}, pa, slogHandler)
	if err != nil {
		panic(err)
	}

	datadir := "/home/lucas/dev/brage/data"
	scfg, err := config.GetServerConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	scfg.DBPath = filepath.Join(datadir, "config/data.db")
	dbSqlite, err := sqlx.Open("sqlite3", scfg.DBPath)
	if err != nil {
		panic(err)
	}
	defer dbSqlite.Close()

	db, err := database.New(dbSqlite)
	if err != nil {
		panic(err)
	}

	t, err := db.GetEnhancedTracksFromAlbumID(context.Background(), "96ad8a93-36fc-41c9-90da-5d0371b2a0da")
	if err != nil {
		panic(err)
	}

	err = ap.LoadAndStartTracks(context.Background(), t, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(ap.CurrentTrack())
	time.Sleep(10 * time.Second)
}
