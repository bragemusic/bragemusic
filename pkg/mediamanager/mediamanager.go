package mediamanager

import (
	"log/slog"

	"github.com/bragemusic/core/pkg/database"
)

type MediaManager struct {
	log *slog.Logger
	db  database.DatabaseFace
}

func New(slogHandler slog.Handler, db database.DatabaseFace) MediaManager {
	return MediaManager{
		log: slog.New(slogHandler).With("service", "media-manager"),
		db:  db,
	}
}
