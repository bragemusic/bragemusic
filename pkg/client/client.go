package client

import (
	"context"
	"log/slog"

	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/serverclient"
)

type Config struct{}

type Client struct {
	sc     *serverclient.ServerClient
	mm     *mediamanager.MediaManager
	config Config
	log    *slog.Logger
}

func (c Client) Sync(ctx context.Context) error {
	return nil
}

func New(sc *serverclient.ServerClient, mm *mediamanager.MediaManager, config Config, slogHandler slog.Handler) Client {
	return Client{
		sc:     sc,
		mm:     mm,
		config: config,
		log:    slog.New(slogHandler).With("service", "client"),
	}
}
