package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddPlayHistory(ctx context.Context, trackID, userID uuid.UUID) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := uid.String()
	playedAt := time.Now()

	ph := types.PlayHistory{
		ID:       id,
		UserID:   userID.String(),
		TrackID:  trackID.String(),
		PlayedAt: playedAt,
	}

	return d.AddPlayHistoryStruct(ctx, ph)
}

func (d Database) AddPlayHistoryStruct(ctx context.Context, ph types.PlayHistory) (string, error) {
	query := `
        INSERT INTO play_history (
            id, user_id, track_id, played_at
        )
        VALUES (?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		ph.ID,
		ph.UserID,
		ph.TrackID,
		ph.PlayedAt,
	)
	if err != nil {
		return "", err
	}

	return ph.ID, nil
}

func (d Database) ListUpdatedPlayHistory(ctx context.Context, since time.Time) (updatedItems []types.PlayHistory, err error) {
	query := `
        SELECT *
        FROM play_history
        WHERE
          played_at > ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &updatedItems, query, since)
	if err != nil {
		return nil, err
	}

	return
}
