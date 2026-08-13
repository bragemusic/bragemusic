package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddTrackArtist(ctx context.Context, ta types.TrackArtist, userID uuid.UUID) (uuid.UUID, error) {
	if ta.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		ta.ID = uid
	}

	now := time.Now()

	if ta.CreatedAt.IsZero() {
		ta.CreatedAt = now
	}

	if ta.UpdatedAt.IsZero() {
		ta.UpdatedAt = now
	}

	query := `
        INSERT INTO track_artists (
            id, track_id , artist_id, role,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		ta.ID,
		ta.TrackID,
		ta.ArtistID,
		ta.Role,
		ta.CreatedAt,
		ta.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, ta.ID, types.EntityEventCreate, types.EntityTrackArtist, userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return ta.ID, nil
}

func (d Database) DeleteTrackArtist(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM track_artists
		WHERE
			id = ?;
	`

	_, err := d.ext.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityTrackArtist, userID)
}

func (d Database) GetTrackArtistByID(ctx context.Context, id uuid.UUID) (trackArtist types.TrackArtist, err error) {
	query := `
		SELECT *
		FROM track_artists
		WHERE
			id = ?;
    `
	err = sqlx.GetContext(ctx, d.ext, &trackArtist, query, id)
	if err != nil {
		return types.TrackArtist{}, err
	}

	return
}

func (d Database) ListTrackArtistsByTrackID(ctx context.Context, trackID uuid.UUID) (trackArtists []types.TrackArtist, err error) {
	query := `
		SELECT *
		FROM track_artists
		WHERE
			track_id = ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &trackArtists, query, trackID)
	if err != nil {
		return nil, err
	}

	return
}
