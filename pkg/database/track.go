package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddTracks(ctx context.Context, tracks []types.Track) (ids []string, err error) {
	for _, track := range tracks {
		trackExistsMbID := false
		if track.MusicBrainzID != nil {
			trackExistsMbID, err = d.TrackExistsByMbID(ctx, *track.MusicBrainzID)
			if err != nil {
				return nil, err
			}
		}

		trackExists, err := d.TrackExistsByNameAndAlbumID(ctx, track.Title, *track.AlbumID)
		if err != nil {
			return nil, err
		}

		if trackExistsMbID || trackExists {
			if err = d.UpdateTrackFromMbID(ctx, track); err != nil {
				return nil, err
			}
			ids = append(ids, track.ID)
		} else {
			id, err := d.AddTrack(ctx, track)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}

	}

	return ids, nil
}

func (d Database) AddTrack(ctx context.Context, t types.Track) (string, error) {
	if t.ID == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		t.ID = uid.String()
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := `
        INSERT INTO tracks (
            id, title, album_id, musicbrainz_id, track_artist, track_number, disc_number,
            genre, year, composer, comment, duration_ms, bitrate, sample_rate,
            file_path, file_size, mime_type,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		t.ID,
		t.Title,
		t.AlbumID,
		t.MusicBrainzID,
		t.TrackArtist,
		t.TrackNumber,
		t.DiscNumber,
		t.Genre,
		t.Year,
		t.Composer,
		t.Comment,
		t.DurationMS,
		t.Bitrate,
		t.SampleRate,
		t.FilePath,
		t.FileSize,
		t.MimeType,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return "", err
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
	query := `
        UPDATE tracks SET
            title = :title,
            album_id = :album_id,
            musicbrainz_id = :musicbrainz_id,
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

func (d Database) GetTrackFromID(ctx context.Context, ID string) (track types.Track, err error) {
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

func (d Database) GetTracksFromAlbumID(ctx context.Context, albumID string) (tracks []types.Track, err error) {
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

func (d Database) GetEnhancedTracksFromAlbumID(ctx context.Context, albumID string) (tracks []types.TrackEnhanced, err error) {
	query := `
        SELECT
            t.*,
            al.name  AS album_name,
            ar.id    AS artist_id,
            ar.name  AS artist_name
        FROM tracks t
        JOIN albums al ON t.album_id = al.id
        JOIN artists ar ON al.artist_id = ar.id
        WHERE t.album_id = ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &tracks, query, albumID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) GetTrackFromName(ctx context.Context, albumID string, trackName string) (track types.Track, err error) {
	query := `
        SELECT *
        FROM tracks
        WHERE album_id = ?
        AND normalize(title) = normalize(?)
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &track, query, albumID, trackName)
	if err != nil {
		return types.Track{}, err
	}

	return
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
