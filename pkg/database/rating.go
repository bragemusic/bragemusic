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

func (d Database) AddRating(ctx context.Context, r types.Rating) (uuid.UUID, error) {
	if r.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		r.ID = uid
	}

	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	query := `
        INSERT INTO ratings (
            id, track_id, rating, owner,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		r.ID,
		r.TrackID,
		r.Rating,
		r.Owner,
		r.CreatedAt,
		r.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return r.ID, nil
}

func (d Database) GetRatingID(ctx context.Context, trackID, userID uuid.UUID) (id uuid.UUID, found bool, err error) {
	query := `
        SELECT id
        FROM ratings
        WHERE track_id = ?
          AND owner = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &id, query, trackID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, false, nil
		}
		return uuid.Nil, false, err
	}

	return id, true, nil
}

func (d Database) UpdateRating(ctx context.Context, id uuid.UUID, rating int) error {
	query := `
        UPDATE ratings
        SET
            rating = ?,
            updated_at = ?
        WHERE id = ?;
    `

	_, err := d.ext.ExecContext(ctx, query, rating, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
