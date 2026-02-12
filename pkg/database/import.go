package database

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type ImportFace interface {
	Begin(ctx context.Context) (DatabaseFace, error)
	driver.Tx

	AddImport(ctx context.Context, i types.Import) (uuid.UUID, error)
}

func (d Database) AddImport(ctx context.Context, i types.Import) (uuid.UUID, error) {
	if i.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		i.ID = uid
	}

	now := time.Now()
	i.CreatedAt = now
	i.UpdatedAt = now

	query := `
        INSERT INTO imports (
            id, musicbrainz_id, owner, filename, type, state,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		i.ID,
		i.MusicBrainzID,
		i.Owner,
		i.Filename,
		i.Type,
		i.State,
		i.CreatedAt,
		i.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return i.ID, nil
}
