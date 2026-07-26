package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddAlbumArtist(ctx context.Context, aa types.AlbumArtist, userID uuid.UUID) (uuid.UUID, error) {
	if aa.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		aa.ID = uid
	}

	now := time.Now()

	if aa.CreatedAt.IsZero() {
		aa.CreatedAt = now
	}

	if aa.UpdatedAt.IsZero() {
		aa.UpdatedAt = now
	}

	query := `
        INSERT INTO album_artists (
            id, album_id, artist_id, role, position,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		aa.ID,
		aa.AlbumID,
		aa.ArtistID,
		aa.Role,
		aa.Position,
		aa.CreatedAt,
		aa.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, aa.ID, types.EntityEventCreate, types.EntityAlbumArtist, userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return aa.ID, nil
}

func (d Database) AlbumArtistExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
        SELECT COUNT(1)
        FROM album_artists
        WHERE id = ?;
	`

	var exists bool
	err := sqlx.GetContext(ctx, d.ext, &exists, query, id)
	if err != nil {
		return false, err
	}

	return exists, nil
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
	if len(tracks) == 0 {
		return nil
	}

	trackIDs := make([]string, 0, len(tracks))
	trackIndex := make(map[string]*types.TrackDetailed)

	for i := range tracks {
		trackIDs = append(trackIDs, tracks[i].ID.String())
		trackIndex[tracks[i].ID.String()] = &tracks[i]
	}

	query := `
		SELECT
			at.track_id,
            aa.role,
			ar.id   AS artist_id,
			ar.name AS artist_name,
            ar.sort_name AS sort_name
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
			trackID        string
			role           types.ArtistRole
			artistID       uuid.UUID
			artistName     string
			artistSortName string
		)

		if err := rows.Scan(&trackID, &role, &artistID, &artistName, &artistSortName); err != nil {
			return err
		}

		if t := trackIndex[trackID]; t != nil {
			t.Artists = append(t.Artists, types.ArtistMinimal{
				ID:       artistID,
				Name:     artistName,
				SortName: artistSortName,
				Role:     role,
			})
		}
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	query2 := `
		SELECT
			at.track_id,
            ta.role,
			ar.id   AS artist_id,
			ar.name AS artist_name,
            ar.sort_name AS sort_name
		FROM album_tracks at
		JOIN track_artists ta ON ta.track_id = at.track_id
		JOIN artists ar ON ar.id = ta.artist_id
		WHERE at.track_id IN (?);
	`

	query2, args, err = sqlx.In(query2, trackIDs)
	if err != nil {
		return err
	}
	query2 = d.ext.Rebind(query2)

	rows2, err := d.ext.QueryxContext(ctx, query2, args...)
	if err != nil {
		return err
	}
	defer rows2.Close()

	for rows2.Next() {
		var (
			trackID        string
			role           types.ArtistRole
			artistID       uuid.UUID
			artistName     string
			artistSortName string
		)

		if err := rows2.Scan(&trackID, &role, &artistID, &artistName, &artistSortName); err != nil {
			return err
		}

		if t := trackIndex[trackID]; t != nil {
			t.Artists = append(t.Artists, types.ArtistMinimal{
				ID:       artistID,
				Name:     artistName,
				SortName: artistSortName,
				Role:     role,
			})
		}
	}

	return rows2.Err()
}

func (d Database) ListUpdatedAlbumArtists(ctx context.Context, since time.Time) (albumArtists []uuid.UUID, err error) {
	query := `
		SELECT
			id
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

func (d Database) GetAlbumArtistByID(ctx context.Context, id uuid.UUID) (albumArtist types.AlbumArtist, err error) {
	query := `
		SELECT *
		FROM album_artists
		WHERE
			id = ?;
    `
	err = sqlx.GetContext(ctx, d.ext, &albumArtist, query, id)
	if err != nil {
		return types.AlbumArtist{}, err
	}

	return
}

func (d Database) UpdateAlbumArtist(ctx context.Context, aa types.AlbumArtist, userID uuid.UUID) error {
	aa.UpdatedAt = time.Now()
	query := `
        UPDATE album_artists SET
            position = :position,
            updated_at = :updated_at
        WHERE
            album_id = :album_id
            AND artist_id = :artist_id
            AND role = :role;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, aa)
	if err != nil {
		return err
	}

	err = d.addEntityEvent(ctx, aa.ID, types.EntityEventCreate, types.EntityAlbumArtist, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) ListAlbumArtistsByAlbumID(ctx context.Context, albumID uuid.UUID) (albumArtists []types.AlbumArtist, err error) {
	query := `
		SELECT *
		FROM album_artists
		WHERE
			album_id = ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &albumArtists, query, albumID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) DeleteAlbumArtist(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM album_artists
		WHERE
			id = ?;
	`

	_, err := d.ext.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityAlbumArtist, userID)
}

func (d Database) CountAlbumsByArtist(ctx context.Context, artistID uuid.UUID) (int, error) {
	const query = `
        SELECT COUNT(1)
        FROM album_artists
        WHERE artist_id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, artistID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
