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

	if r.CreatedAt.IsZero() {
		now := time.Now()
		r.CreatedAt = now
		r.UpdatedAt = now
	}

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

	err = d.addEntityEvent(ctx, r.ID, types.EntityEventCreate, types.EntityRating)
	if err != nil {
		return uuid.UUID{}, err
	}

	return r.ID, nil
}

func (d Database) GetRating(ctx context.Context, id uuid.UUID) (rating types.Rating, err error) {
	query := `
        SELECT *
        FROM ratings
        WHERE id = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &rating, query, id)
	if err != nil {
		return types.Rating{}, err
	}

	return rating, err
}

func (d Database) GetRatingID(ctx context.Context, trackID, userID uuid.UUID) (uuid.UUID, bool, error) {
	query := `
        SELECT id
        FROM ratings
        WHERE track_id = ?
          AND owner = ?
        LIMIT 1;
    `

	var idStr string
	err := sqlx.GetContext(ctx, d.ext, &idStr, query, trackID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, false, nil
		}
		return uuid.Nil, false, err
	}

	id, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, false, err
	}

	return id, true, nil
}

func (d Database) GetTrackRatings(ctx context.Context, trackID uuid.UUID) (ratings []types.Rating, err error) {
	query := `
        SELECT *
        FROM ratings
        WHERE track_id = ?
        ORDER BY updated_at;
    `

	err = sqlx.SelectContext(ctx, d.ext, &ratings, query, trackID)
	if err != nil {
		return nil, err
	}

	return ratings, nil
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

	err = d.addEntityEvent(ctx, id, types.EntityEventUpdate, types.EntityRating)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) attachTrackRatings(ctx context.Context, tracks []types.TrackDetailed, userID uuid.UUID) error {
	if len(tracks) == 0 {
		return nil
	}

	trackIDs := make([]string, 0, len(tracks))
	trackIndex := make(map[string]*types.TrackDetailed)

	for i := range tracks {
		trackIDs = append(trackIDs, tracks[i].ID)
		trackIndex[tracks[i].ID] = &tracks[i]
	}

	query := `
        SELECT
            r.track_id,
            AVG(r.rating) AS mean_rating,
            MAX(CASE WHEN r.owner = ? THEN r.rating END) AS user_rating
        FROM ratings r
        WHERE r.track_id IN (?)
        GROUP BY r.track_id;
    `

	query, args, err := sqlx.In(query, userID, trackIDs)
	if err != nil {
		return err
	}

	query = d.ext.Rebind(query)

	rows, err := d.ext.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			trackID    string
			meanRating *float64
			userRating *int
		)

		if err := rows.Scan(&trackID, &meanRating, &userRating); err != nil {
			return err
		}

		if t := trackIndex[trackID]; t != nil {
			t.Rating = meanRating
			t.UserRating = userRating
		}
	}

	return rows.Err()
}
