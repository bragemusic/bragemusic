package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddAlbumArtist(ctx context.Context, aa types.AlbumArtist) error {
	now := time.Now()
	aa.CreatedAt = now
	aa.UpdatedAt = now

	query := `
        INSERT INTO album_artists (
            album_id, artist_id, role, position,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		aa.AlbumID,
		aa.ArtistID,
		aa.Role,
		aa.Position,
		aa.CreatedAt,
		aa.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) AlbumArtistExists(ctx context.Context, albumID uuid.UUID, artistID uuid.UUID, role types.ArtistRole) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM album_artists
			WHERE album_id = ?
			  AND artist_id = ?
              AND role = ?
		);
	`

	var exists bool
	err := sqlx.GetContext(ctx, d.ext, &exists, query, albumID, artistID, role)
	if err != nil {
		return false, err
	}

	return exists, nil
}
