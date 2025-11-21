package metasyncer

import (
	"context"
	"log/slog"

	"github.com/bragemusic/core/pkg/database"
)

type MetaSyncer struct {
	db  database.DatabaseFace
	log *slog.Logger
}

func (m MetaSyncer) Sync(ctx context.Context) {

}

func New(db database.DatabaseFace, slogHandler slog.Handler) MetaSyncer {
	return MetaSyncer{
		db:  db,
		log: slog.New(slogHandler).With("service", "meta-syncer"),
	}
}
