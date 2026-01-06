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

	if aa.CreatedAt.IsZero() {
		aa.CreatedAt = now
	}

	if aa.UpdatedAt.IsZero() {
		aa.UpdatedAt = now
	}

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

func (d Database) attachTrackArtists(ctx context.Context, tracks []types.TrackDetailed) error {
	trackIDs := make([]string, 0, len(tracks))
	trackIndex := make(map[string]*types.TrackDetailed)

	for i := range tracks {
		trackIDs = append(trackIDs, tracks[i].ID)
		trackIndex[tracks[i].ID] = &tracks[i]
	}

	query := `
		SELECT
			at.track_id,
			ar.id   AS artist_id,
			ar.name AS artist_name
		FROM album_tracks at
		JOIN album_artists aa ON aa.album_id = at.album_id
		JOIN artists ar ON ar.id = aa.artist_id
		WHERE at.track_id IN (?);
	`

	query, args, err := sqlx.In(query, trackIDs)
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
			artistID   string
			artistName string
		)

		if err := rows.Scan(&trackID, &artistID, &artistName); err != nil {
			return err
		}

		if t := trackIndex[trackID]; t != nil {
			t.ArtistIDs = append(t.ArtistIDs, artistID)
			t.ArtistNames = append(t.ArtistNames, artistName)
		}
	}

	return rows.Err()
}

func (d Database) ListUpdatedAlbumArtists(ctx context.Context, since time.Time) (albumArtists []types.AlbumArtistKey, err error) {
	query := `
		SELECT
			album_id,
			artist_id,
			role
		FROM album_artists
		WHERE
			created_at > ?
			OR updated_at > ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &albumArtists, query, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) GetAlbumArtist(ctx context.Context, albumID, artistID uuid.UUID, role types.ArtistRole) (albumArtist types.AlbumArtist, err error) {
	query := `
		SELECT *
		FROM album_artists
		WHERE
			album_id = ?
			AND artist_id = ?
			AND role = ?;
    `
	err = sqlx.GetContext(ctx, d.ext, &albumArtist, query, albumID, artistID, role)
	if err != nil {
		return types.AlbumArtist{}, err
	}

	return
}

func (d Database) UpdateAlbumArtist(ctx context.Context, aa types.AlbumArtist) error {
	query := `
        UPDATE album_artists SET
            position = :position,
        WHERE
            album_id = :album_id
            AND artist_id = :artist_id
            AND role = :role;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, aa)
	return err
}
