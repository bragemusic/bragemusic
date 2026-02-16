package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddPlaylist(ctx context.Context, p types.Playlist, userID uuid.UUID) (uuid.UUID, error) {
	if p.Owner == uuid.Nil {
		return uuid.Nil, ErrNoUser
	}

	if p.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		p.ID = uid
	}

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	const query = `
		INSERT INTO playlists (
			id, name, description, owner, public,
            created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		p.ID,
		p.Name,
		p.Description,
		p.Owner,
		p.Public,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, p.ID, types.EntityEventCreate, types.EntityPlaylist, userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return p.ID, nil
}

func (d Database) AddPlaylistTrack(ctx context.Context, p types.PlaylistTrack, userID uuid.UUID) (uuid.UUID, error) {
	if p.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		p.ID = uid
	}

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	const query = `
		INSERT INTO playlist_tracks (
			id, playlist_id, album_track_id,
            created_at, updated_at
		) VALUES (?, ?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		p.ID,
		p.PlaylistID,
		p.AlbumTrackID,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, p.ID, types.EntityEventCreate, types.EntityPlaylistTrack, userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return p.ID, nil
}

func (d Database) CountPlaylists(ctx context.Context, userID uuid.UUID) (int, error) {
	const query = `
        SELECT COUNT(1)
        FROM playlists
        WHERE owner = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d Database) CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error) {
	const query = `
        SELECT COUNT(1)
        FROM playlist_tracks
        WHERE playlist_id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, playlistID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d Database) DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM playlists
		WHERE
			id = ?
          AND
            owner = ?;
	`

	_, err := d.ext.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	return d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityPlaylist, userID)
}

func (d Database) DeletePlaylistTrack(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM playlist_tracks
		WHERE
			id = ?
        ;
	`

	_, err := d.ext.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityPlaylistTrack, userID)
}

func (d Database) GetPlaylist(ctx context.Context, ID, userID uuid.UUID) (plist types.Playlist, err error) {
	query := `
        SELECT *
        FROM playlists
        WHERE id = ?
          AND (owner = ? OR public = 1)
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &plist, query, ID, userID)
	if err != nil {
		return types.Playlist{}, err
	}

	return
}

func (d Database) GetPlaylistTrack(ctx context.Context, id uuid.UUID) (plistTrack types.PlaylistTrack, err error) {
	query := `
        SELECT *
        FROM playlist_tracks
        WHERE id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &plistTrack, query, id)
	if err != nil {
		return types.PlaylistTrack{}, err
	}

	return
}

func (d Database) GetPlaylistTrackByPlaylistAndAlbumTrack(ctx context.Context, playlistID, albumTrackID uuid.UUID) (plistTrack types.PlaylistTrack, err error) {
	query := `
        SELECT *
        FROM playlist_tracks
        WHERE playlist_id = ?
          AND album_track_id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &plistTrack, query, playlistID, albumTrackID)
	if err != nil {
		return types.PlaylistTrack{}, err
	}

	return
}

func (d Database) ListPlaylists(ctx context.Context, userID uuid.UUID, includePublic bool, sortBy SortBy, sortOrder SortOrder) (playlists []types.Playlist, err error) {
	sortByStr := ""

	switch sortBy {
	case SortByDate:
		sortByStr = "created_at"
	case SortByName:
		sortByStr = "name"
	case SortByPlayCount:
		sortByStr = "name"
	}

	publicStr := ""
	if includePublic {
		publicStr = "OR public = 1"
	}

	query := fmt.Sprintf(`
        SELECT *
        FROM playlists
        WHERE
          owner = ?
          %s
        ORDER BY %s %s
        ;
    `, publicStr, sortByStr, sortOrder)
	err = sqlx.SelectContext(ctx, d.ext, &playlists, query, userID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
	tracksQuery := `
        SELECT
            t.id,
            t.title,
            at.album_id,
            al.name AS album_name,
            t.musicbrainz_id,
            at.track_number,
            at.disc_number,
            t.genre,
            t.comment,
            t.created_at,
            t.updated_at,
            COALESCE(tp.play_count, 0) AS play_count,
            pt.id as context_id
        FROM playlist_tracks pt
        JOIN album_tracks at ON at.id = pt.album_track_id
        JOIN tracks t ON t.id = at.track_id
        JOIN albums al ON al.id = at.album_id
        LEFT JOIN (
            SELECT track_id, COUNT(*) AS play_count
            FROM play_history
            GROUP BY track_id
        ) tp ON tp.track_id = t.id
        WHERE pt.playlist_id = ?
        ORDER BY
            pt.created_at;
	`

	if err := sqlx.SelectContext(ctx, d.ext, &tracks, tracksQuery, playlistID); err != nil {
		return nil, err
	}

	if err := d.attachTrackArtists(ctx, tracks); err != nil {
		return nil, err
	}

	if err := d.attachMediaFiles(ctx, tracks); err != nil {
		return nil, err
	}

	if err := d.attachTrackRatings(ctx, tracks, userID); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (d Database) ListUpdatedPlaylists(ctx context.Context, since time.Time, userID uuid.UUID) (plists []uuid.UUID, err error) {
	query := `
		SELECT
			id
		FROM playlists
		WHERE
            owner = ?
          AND
			(created_at > ?
			OR updated_at > ?);
    `
	err = sqlx.SelectContext(ctx, d.ext, &plists, query, userID, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListUpdatedPlaylistTracks(ctx context.Context, since time.Time, userID uuid.UUID) (plists []uuid.UUID, err error) {
	query := `
       SELECT
           pt.id
       FROM playlist_tracks pt
       JOIN playlists p ON p.id = pt.playlist_id
       WHERE
           p.owner = ?
           AND (
               pt.created_at > ?
               OR pt.updated_at > ?
           );
    `
	err = sqlx.SelectContext(ctx, d.ext, &plists, query, userID, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) PlaylistExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
        SELECT COUNT(1)
        FROM playlists
        WHERE id = ?;
	`

	var exists bool
	err := sqlx.GetContext(ctx, d.ext, &exists, query, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (d Database) PlaylistTrackExists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
        SELECT COUNT(1)
        FROM playlist_tracks
        WHERE id = ?;
	`

	var exists bool
	err := sqlx.GetContext(ctx, d.ext, &exists, query, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (d Database) UpdatePlaylist(ctx context.Context, plist types.Playlist, userID uuid.UUID) error {
	plist.UpdatedAt = time.Now()
	query := `
        UPDATE playlists SET
            name = :name,
            description = :description,
            updated_at = :updated_at
        WHERE
            id = :id
            AND owner = :owner;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, plist)
	if err != nil {
		return err
	}

	err = d.addEntityEvent(ctx, plist.ID, types.EntityEventUpdate, types.EntityPlaylist, userID)
	if err != nil {
		return err
	}

	return nil
}
