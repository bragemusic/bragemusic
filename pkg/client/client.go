package client

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/migrations"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	ServerBaseURL string
	MusicDirPath  string
	ConfigPath    string
	ImagePath     string
}

type Client struct {
	sc      *serverclient.ServerClient
	mm      *mediamanager.MediaManager
	sy      *syncer.Syncer
	config  Config
	log     *slog.Logger
	dbClose func() error
}

func (c Client) Sync(ctx context.Context) error {
	return c.sy.Sync(ctx)
}

func (c *Client) Close() error {
	return c.dbClose()
}

func NewSyncer(config Config, slogHandler slog.Handler) (c Client, err error) {
	dbPath := filepath.Join(config.ConfigPath, "data.db")
	if err = migrations.Migrate(context.Background(), dbPath, slogHandler); err != nil {
		return Client{}, err
	}

	dbSqlite, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return Client{}, err
	}

	db, err := database.New(dbSqlite)
	if err != nil {
		return Client{}, err
	}

	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	mm := mediamanager.New(slogHandler, &db)
	sy := syncer.New(&sc, &db, config.MusicDirPath, slogHandler)

	return Client{
		sc:      &sc,
		mm:      &mm,
		sy:      &sy,
		config:  config,
		log:     slog.New(slogHandler).With("service", "client"),
		dbClose: dbSqlite.Close,
	}, nil
}
