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

func (d Database) AddTrackAnalysis(ctx context.Context, ta types.TrackAnalysis) error {
	if ta.ID == uuid.Nil {
		return errors.New("track_analysis must have an ID")
	}

	now := time.Now()
	ta.CreatedAt = now
	ta.UpdatedAt = now

	query := `
        INSERT INTO track_analyses (
            id, bpm, key, key_scale, key_confidence, loudness, energy, danceability, mood_happy, mood_sad, mood_aggresive, mood_calm,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		ta.ID,
		ta.BPM,
		ta.Key,
		ta.KeyScale,
		ta.KeyConfidence,
		ta.Loudness,
		ta.Energy,
		ta.Danceability,
		ta.MoodHappy,
		ta.MoodSad,
		ta.MoodAggresive,
		ta.MoodCalm,
		ta.CreatedAt,
		ta.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) GetUnanalysedTrack(ctx context.Context) (id uuid.UUID, found bool, err error) {
	query := `
        SELECT t.id
        FROM tracks t
        LEFT JOIN track_analyses ta ON ta.id = t.id
        WHERE ta.id IS NULL
        AND t.media_file IS NOT NULL
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &id, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, false, nil
		}
		return uuid.Nil, false, err
	}

	return id, true, nil
}
