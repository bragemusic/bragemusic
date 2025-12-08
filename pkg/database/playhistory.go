package database

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

func (d Database) AddPlayHistory(ctx context.Context, trackID, userID uuid.UUID) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := uid.String()
	playedAt := time.Now()

	query := `
        INSERT INTO play_history (
            id, user_id, track_id, played_at
        )
        VALUES (?, ?, ?, ?);
    `

	_, err = d.ext.ExecContext(
		ctx,
		query,
		id,
		userID.String(),
		trackID.String(),
		playedAt,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}
