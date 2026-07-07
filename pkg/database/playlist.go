package database

import (
	"context"
	"database/sql"
	"errors"
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

	p.Type = types.PlaylistTypeStandard

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	const query = `
		INSERT INTO playlists (
			id, name, description, owner, public, type,
            created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		p.ID,
		p.Name,
		p.Description,
		p.Owner,
		p.Public,
		p.Type,
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

func (d Database) AddSmartPlaylist(ctx context.Context, p types.SmartPlaylist, userID uuid.UUID) (uuid.UUID, error) {
	if p.Owner == uuid.Nil {
		return uuid.Nil, ErrNoUser
	}

	p.Type = types.PlaylistTypeSmart

	playlistID, err := d.AddPlaylist(ctx, p.Playlist, userID)
	if err != nil {
		return uuid.Nil, err
	}

	p.Content.PlaylistID = playlistID

	c := p.Content

	if c.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		c.ID = uid
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	const query = `
		INSERT INTO smart_playlist_contents (
			id, playlist_id, bpm_upper, bpm_lower, mood_happy, mood_sad, mood_aggressive, mood_calm,
            created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err = d.ext.ExecContext(
		ctx,
		query,
		c.ID,
		c.PlaylistID,
		c.BpmUpper,
		c.BpmLower,
		c.MoodHappy,
		c.MoodSad,
		c.MoodAggressive,
		c.MoodCalm,
		c.CreatedAt,
		c.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, c.ID, types.EntityEventCreate, types.EntitySmartPlaylistContent, userID)
	if err != nil {
		return uuid.Nil, err
	}

	if c.Artists != nil {
		for _, a := range *c.Artists {
			_, err := d.addSmartPlaylistArtist(ctx, c.ID, a, userID)
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	return p.ID, nil
}

func (d Database) addSmartPlaylistArtist(ctx context.Context, parentID, artistID, userID uuid.UUID) (uuid.UUID, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	now := time.Now()

	const query = `
		INSERT INTO smart_playlist_artists (
			id, parent_id, artist_id,
            created_at, updated_at
		) VALUES (?, ?, ?, ?, ?);
	`

	_, err = d.ext.ExecContext(
		ctx,
		query,
		uid,
		parentID,
		artistID,
		now,
		now,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, uid, types.EntityEventCreate, types.EntitySmartPlaylistArtist, userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uid, nil
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

func (d Database) listPlaylistStandardTracks(ctx context.Context, playlistID, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
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

	return tracks, nil
}

func (d Database) listPlaylistSmartTracks(ctx context.Context, playlistID, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
	contentQuery := `
		SELECT * FROM smart_playlist_contents
		WHERE playlist_id = ?
		LIMIT 1;
`

	content := types.SmartPlaylistContent{}
	err = sqlx.GetContext(ctx, d.ext, &content, contentQuery, playlistID)
	if err != nil {
		return nil, err
	}

	artistQuery := `
		SELECT * FROM smart_playlist_artists
		WHERE parent_id = ?;
`

	artists := []types.SmartPlaylistArtist{}
	if err := sqlx.SelectContext(ctx, d.ext, &artists, artistQuery, content.ID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if len(artists) > 0 {
		ids := []uuid.UUID{}
		for _, a := range artists {
			ids = append(ids, a.ArtistID)
		}
		content.Artists = &ids
	}

	filter := content.TrackFilter()

	tracksNew, totalPages, _, err := d.ListTracksWithFilters(ctx, filter, 1, 100)
	if err != nil {
		return nil, err
	}

	page := 2

	for page <= totalPages {
		tr, _, _, err := d.ListTracksWithFilters(ctx, filter, page, 100)
		if err != nil {
			return nil, err
		}

		tracksNew = append(tracksNew, tr...)
		page += 1
	}

	for _, t := range tracksNew {
		tracks = append(tracks, types.TrackDetailed{
			ID:            t.Track.ID,
			Title:         t.Track.Title,
			AlbumID:       t.Album.ID.String(),
			AlbumName:     t.Album.Name,
			MusicBrainzID: t.Track.MusicBrainzID,
			TrackNumber:   t.AlbumTrack.TrackNumber,
			DiscNumber:    t.AlbumTrack.DiscNumber,
			Genre:         nil,
			Comment:       nil,
			MediaFile:     t.Mediafile,
			PlayCount:     0,
			ContextID:     nil,
			Rating:        nil,
			UserRating:    nil,
			Liked:         false,
			CreatedAt:     t.Track.CreatedAt,
			UpdatedAt:     t.Track.UpdatedAt,
		})
	}

	return tracks, nil
}

func (d Database) ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
	playlist, err := d.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return nil, err
	}

	switch playlist.Type {
	case types.PlaylistTypeStandard:
		tracks, err = d.listPlaylistStandardTracks(ctx, playlistID, userID)
		if err != nil {
			return nil, err
		}
	case types.PlaylistTypeSmart:
		tracks, err = d.listPlaylistSmartTracks(ctx, playlistID, userID)
		if err != nil {
			return nil, err
		}
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

	if err := d.attachTrackLike(ctx, tracks, userID); err != nil {
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
