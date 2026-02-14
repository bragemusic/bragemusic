package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddTrack(ctx context.Context, t types.Track) (uuid.UUID, error) {
	if t.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		t.ID = uid
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := `
        INSERT INTO tracks (
            id, title, musicbrainz_id,
            genre, comment, media_file,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		t.ID,
		t.Title,
		t.MusicBrainzID,
		t.Genre,
		t.Comment,
		t.MediaFile,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return t.ID, nil
}

func (d Database) TrackExists(ctx context.Context, ID string) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM tracks
        WHERE id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, ID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) TrackExistsByNameAndAlbumID(ctx context.Context, title, albumID string) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM tracks
        WHERE title = ? AND album_id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, title, albumID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) TrackExistsByMbID(ctx context.Context, trackMbID string) (bool, error) {
	query := `
        SELECT COUNT(1)
        FROM tracks
        WHERE musicbrainz_id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, trackMbID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) UpdateTrack(ctx context.Context, t types.Track) error {
	t.UpdatedAt = time.Now()
	query := `
        UPDATE tracks SET
            title = :title,
            musicbrainz_id = :musicbrainz_id,
            genre = :genre,
            comment = :comment,
            media_file = :media_file,
            updated_at = :updated_at
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, t)
	return err
}

func (d Database) UpdateTrackFromMbID(ctx context.Context, t types.Track) error {
	query := `
        UPDATE tracks SET
            title = :title,
            album_id = :album_id,
            track_artist = :track_artist,
            track_number = :track_number,
            disc_number = :disc_number,
            genre = :genre,
            year = :year,
            composer = :composer,
            comment = :comment,
            duration_ms = :duration_ms,
            bitrate = :bitrate,
            sample_rate = :sample_rate,
            file_path = :file_path,
            file_size = :file_size,
            mime_type = :mime_type
        WHERE musicbrainz_id = :musicbrainz_id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, t)
	return err
}

func (d Database) GetTrackDetailed(ctx context.Context, trackID, albumID uuid.UUID) (track types.TrackDetailed, err error) {
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
			COALESCE(tp.play_count, 0) AS play_count
		FROM album_tracks at
		JOIN tracks t ON t.id = at.track_id
		JOIN albums al ON al.id = at.album_id
		LEFT JOIN (
			SELECT track_id, COUNT(*) AS play_count
			FROM play_history
			GROUP BY track_id
		) tp ON tp.track_id = t.id
		WHERE at.album_id = ?
        AND t.id = ?
		;
	`

	if err := sqlx.GetContext(ctx, d.ext, &track, tracksQuery, albumID, trackID); err != nil {
		return types.TrackDetailed{}, err
	}

	dummySlice := []types.TrackDetailed{track}

	if err := d.attachTrackArtists(ctx, dummySlice); err != nil {
		return types.TrackDetailed{}, err
	}

	if err := d.attachMediaFiles(ctx, dummySlice); err != nil {
		return types.TrackDetailed{}, err
	}

	return dummySlice[0], nil
}

func (d Database) GetTrackFromID(ctx context.Context, ID uuid.UUID) (track types.Track, err error) {
	query := `
        SELECT *
        FROM tracks
        WHERE id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &track, query, ID)
	if err != nil {
		return types.Track{}, err
	}

	return
}

func (d Database) GetTrackFromMbID(ctx context.Context, mbID string) (track types.Track, err error) {
	query := `
        SELECT *
        FROM tracks
        WHERE musicbrainz_id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &track, query, mbID)
	if err != nil {
		return types.Track{}, err
	}

	return
}

func (d Database) GetTracksFromAlbumID(ctx context.Context, albumID uuid.UUID) (tracks []types.Track, err error) {
	query := `
        SELECT *
        FROM tracks
        WHERE album_id = ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &tracks, query, albumID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) GetTracksDetailedFromArtistID(ctx context.Context, artistID uuid.UUID, sortBy SortBy, sortOrder SortOrder, limit *int, includeMissingFiles bool) (tracks []types.TrackDetailed, err error) {
	orderBy := "t.title"

	switch sortBy {
	case SortByDate:
		orderBy = "t.created_at"
	case SortByName:
		orderBy = "t.title"
	case SortByPlayCount:
		orderBy = "play_count"
	}

	whereExtra := ""
	if !includeMissingFiles {
		whereExtra = "AND t.media_file IS NOT NULL"
	}

	limitClause := ""
	if limit != nil {
		limitClause = fmt.Sprintf("LIMIT %d", *limit)
	}

	query := fmt.Sprintf(`
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
			COALESCE(tp.play_count, 0) AS play_count

		FROM album_artists aa
		JOIN albums al        ON al.id = aa.album_id
		JOIN album_tracks at ON at.album_id = al.id
		JOIN tracks t        ON t.id = at.track_id

		LEFT JOIN (
			SELECT track_id, COUNT(*) AS play_count
			FROM play_history
			GROUP BY track_id
		) tp ON tp.track_id = t.id

		WHERE aa.artist_id = ?
		%s

		ORDER BY %s %s
		%s;
	`, whereExtra, orderBy, sortOrder, limitClause)

	err = sqlx.SelectContext(ctx, d.ext, &tracks, query, artistID)
	if err != nil {
		return nil, err
	}

	// attach optional relations
	if err := d.attachMediaFiles(ctx, tracks); err != nil {
		return nil, err
	}

	if err := d.attachTrackArtists(ctx, tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (d Database) GetTrackFromName(ctx context.Context, albumID uuid.UUID, trackName string) (track types.Track, err error) {
	query := `
		SELECT t.*
		FROM album_tracks at
		JOIN tracks t ON t.id = at.track_id
		WHERE at.album_id = ?
		  AND normalize(t.title) = normalize(?)
		ORDER BY at.disc_number, at.track_number
		LIMIT 1;
	`

	err = sqlx.GetContext(ctx, d.ext, &track, query, albumID.String(), trackName)
	if err != nil {
		return types.Track{}, err
	}

	return track, nil
}

func (d Database) ListTracks(ctx context.Context) (tracks []types.Track, err error) {
	query := `
        SELECT *
        FROM tracks
        ;
    `

	err = sqlx.SelectContext(ctx, d.ext, &tracks, query)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (d Database) ListUpdatedTracks(ctx context.Context, since time.Time) (trackIDs []string, err error) {
	query := `
        SELECT id
        FROM tracks
        WHERE
          created_at > ?
          OR
          updated_at > ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &trackIDs, query, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListAlbumTracksDetailed(ctx context.Context, albumID, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
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
			COALESCE(tp.play_count, 0) AS play_count
		FROM album_tracks at
		JOIN tracks t ON t.id = at.track_id
		JOIN albums al ON al.id = at.album_id
		LEFT JOIN (
			SELECT track_id, COUNT(*) AS play_count
			FROM play_history
			GROUP BY track_id
		) tp ON tp.track_id = t.id
		WHERE at.album_id = ?
		ORDER BY at.disc_number, at.track_number;
	`

	if err := sqlx.SelectContext(ctx, d.ext, &tracks, tracksQuery, albumID); err != nil {
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

func (d Database) CountTracks(ctx context.Context) (int, error) {
	const query = `
        SELECT COUNT(1)
        FROM tracks;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
