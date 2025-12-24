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

	// err = a.RemoveUser(ctx, uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")))
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	return
	// }
	//
	// t, err := a.GenerateToken(ctx, uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")), types.TokenFrontendLong, nil)
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	return
	// }
	// fmt.Println(t)
	token := "brg_v1_TTqTKHxCuJ1y491EZBWGDMFHI4ybwgymEQV6J-MZjew"
	_ = token
	u, err := a.GetUserFromTokenString(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Println(u)

	m := mediamanager.New(slogHandler, db, scfg.Paths.MusicDir)

	s := server.New(slogHandler, &m, scfg)

	logger.Info(fmt.Sprintf("serving on port %d", scfg.Port))
	if err = http.ListenAndServe(fmt.Sprintf(":%d", scfg.Port), s.Handler()); err != nil {
		logger.Error(err.Error())
		return
	}
}
