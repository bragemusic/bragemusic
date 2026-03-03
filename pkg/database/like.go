package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddLike(ctx context.Context, l types.Like) (uuid.UUID, error) {
	if l.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		l.ID = uid
	}

	if l.CreatedAt.IsZero() {
		now := time.Now()
		l.CreatedAt = now
		l.UpdatedAt = now
	}

	query := `
        INSERT INTO likes (
            id, track_id, owner,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		l.ID,
		l.TrackID,
		l.Owner,
		l.CreatedAt,
		l.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, l.ID, types.EntityEventCreate, types.EntityLike, l.Owner)
	if err != nil {
		return uuid.UUID{}, err
	}

	return l.ID, nil
}

func (d Database) DeleteLike(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM likes
		WHERE
			id = ?;
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return err
	}

	err = d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityLike, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) GetLikeID(ctx context.Context, trackID, userID uuid.UUID) (uuid.UUID, error) {
	query := `
        SELECT id
        FROM likes
        WHERE track_id = ?
          AND owner = ?
        LIMIT 1;
    `

	var idStr string
	err := sqlx.GetContext(ctx, d.ext, &idStr, query, trackID, userID)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (d Database) HasLike(ctx context.Context, trackID, userID uuid.UUID) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM likes
        WHERE owner = ?
        AND track_id = ?
;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, userID, trackID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
