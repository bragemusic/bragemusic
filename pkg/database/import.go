package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

type ImportFace interface {
	Begin(ctx context.Context) (DatabaseFace, error)
	driver.Tx

	AddImport(ctx context.Context, i types.Import) (uuid.UUID, error)
	GetUnclaimedImport(ctx context.Context) (i types.Import, found bool, err error)
	SetImportState(ctx context.Context, id uuid.UUID, state types.ImportState) error
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

func (d Database) GetUnclaimedImport(ctx context.Context) (i types.Import, found bool, err error) {
	query := `
        SELECT *
        FROM imports
        WHERE state = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &i, query, types.ImportStateNotStarted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Import{}, false, nil
		}
		return types.Import{}, false, err
	}

	return i, true, nil
}

func (d Database) SetImportState(ctx context.Context, id uuid.UUID, state types.ImportState) error {
	query := `
        UPDATE imports
        SET
            state = ?,
            updated_at = ?
        WHERE id = ?;
    `

	_, err := d.ext.ExecContext(ctx, query, state, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
