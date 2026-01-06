package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddAlbumTrack(ctx context.Context, at types.AlbumTrack) error {
	now := time.Now()
	at.CreatedAt = now
	at.UpdatedAt = now

	query := `
        INSERT INTO album_tracks (
            album_id, track_id, disc_number, track_number,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		at.AlbumID,
		at.TrackID,
		at.DiscNumber,
		at.TrackNumber,
		at.CreatedAt,
		at.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) AlbumTrackExists(ctx context.Context, albumID uuid.UUID, trackID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM album_tracks
			WHERE album_id = ?
			  AND track_id = ?
		);
	`

	var exists bool
	err := sqlx.GetContext(ctx, d.ext, &exists, query, albumID, trackID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (d Database) ListUpdatedAlbumTracks(ctx context.Context, since time.Time) (albumTracks []types.AlbumTrackKey, err error) {
	query := `
		SELECT
			album_id,
            disc_number,
            track_number
		FROM album_tracks
		WHERE
			created_at > ?
			OR updated_at > ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &albumTracks, query, since, since)
	if err != nil {
		return nil, err
	}

	return
}
