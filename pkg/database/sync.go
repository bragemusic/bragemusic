package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("db entry not found")

func (d Database) AddSync(ctx context.Context, s types.DBSyncState) (string, error) {
	if s.ID == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		s.ID = uid.String()
	}

	query := `
        INSERT INTO syncs (
            id, artists_created, artists_updated, albums_created, albums_updated, tracks_created, tracks_updated,
            synced_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		s.ID,
		s.ArtistsCreated,
		s.ArtistsUpdated,
		s.AlbumsCreated,
		s.AlbumsUpdated,
		s.TracksCreated,
		s.TracksUpdated,
		s.SyncedAt,
	)
	if err != nil {
		return "", err
	}

	return s.ID, nil
}

func (d Database) GetLastSync(ctx context.Context) (sync types.DBSyncState, err error) {
	query := `
        SELECT *
        FROM syncs
        ORDER BY synced_at DESC
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &sync, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.DBSyncState{}, ErrNotFound
		}
		return types.DBSyncState{}, err
	}

	return
}

func (d Database) AddSyncItem(ctx context.Context, s types.SyncItem) (uuid.UUID, error) {
	if s.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		s.ID = uid
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	query := `
        INSERT INTO sync_items (
            id, sync_id, item_id, type, state,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		s.ID,
		s.SyncID,
		s.ItemID,
		s.Type,
		s.State,
		s.CreatedAt,
		s.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return s.ID, nil
}
