package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

type trackFilterQuery struct {
	Joins []string
	Where []string
	Args  []any
}

func (d Database) AddTrack(ctx context.Context, t types.Track, userID uuid.UUID) (uuid.UUID, error) {
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

	err = d.addEntityEvent(ctx, t.ID, types.EntityEventCreate, types.EntityTrack, userID)
	if err != nil {
		return uuid.UUID{}, err
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

func (d Database) UpdateTrack(ctx context.Context, t types.Track, userID uuid.UUID) error {
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
	if err != nil {
		return err
	}

	err = d.addEntityEvent(ctx, t.ID, types.EntityEventUpdate, types.EntityTrack, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) GetTrackDetailed(ctx context.Context, trackID, albumID, userID uuid.UUID) (track types.TrackDetailed, err error) {
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

	if err := d.attachTrackRatings(ctx, dummySlice, userID); err != nil {
		return types.TrackDetailed{}, err
	}

	if err := d.attachTrackLike(ctx, dummySlice, userID); err != nil {
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

func (d Database) GetTracksDetailedFromArtistID(ctx context.Context, artistID, userID uuid.UUID, sortBy SortBy, sortOrder SortOrder, limit *int, includeMissingFiles bool) (tracks []types.TrackDetailed, err error) {
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

	if err := d.attachTrackRatings(ctx, tracks, userID); err != nil {
		return nil, err
	}

	if err := d.attachTrackLike(ctx, tracks, userID); err != nil {
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

	if err := d.attachTrackLike(ctx, tracks, userID); err != nil {
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

func (d Database) buildTrackFilterQuery(filter types.TrackFilter) trackFilterQuery {
	q := trackFilterQuery{}

	if filter.BPM != nil {
		q.Where = append(q.Where, "ta.bpm BETWEEN ? AND ?")
		q.Args = append(q.Args, filter.BPM.Lower, filter.BPM.Upper)
	}

	if filter.Mood.Aggressive != nil {
		q.Where = append(q.Where, "ta.mood_aggresive > ?")
		q.Args = append(q.Args, *filter.Mood.Aggressive)
	}

	if filter.Mood.Calm != nil {
		q.Where = append(q.Where, "ta.mood_calm > ?")
		q.Args = append(q.Args, *filter.Mood.Calm)
	}

	if filter.Artists != nil && len(*filter.Artists) > 0 {

		q.Joins = append(q.Joins,
			"JOIN album_artists aa ON aa.album_id = al.id",
		)

		placeholders := strings.Repeat("?,", len(*filter.Artists))
		placeholders = placeholders[:len(placeholders)-1]

		q.Where = append(q.Where,
			"aa.artist_id IN ("+placeholders+")",
		)

		for _, id := range *filter.Artists {
			q.Args = append(q.Args, id)
		}
	}

	return q
}

func (d Database) ListTracksWithFilters(ctx context.Context, filter types.TrackFilter, page, limit int) (results []types.TrackDetailedNew, totalPages, totalItems int, err error) {
	fq := d.buildTrackFilterQuery(filter)

	query1 := `
      SELECT DISTINCT
          t.*,
          al.*,
          at.*,
          mf.*,
          ta.*
      FROM track_analyses ta
      JOIN tracks t ON t.id = ta.id
      JOIN album_tracks at ON at.track_id = t.id
      JOIN albums al ON al.id = at.album_id
      LEFT JOIN media_files mf ON mf.id = t.media_file
`

	if len(fq.Joins) > 0 {
		query1 += "\n" + strings.Join(fq.Joins, "\n")
	}

	if len(fq.Where) > 0 {
		query1 += "\nWHERE " + strings.Join(fq.Where, " AND ")
	}

	query1 += `
      ORDER BY al.sort_name, at.disc_number, at.track_number
      LIMIT ?
      OFFSET ?
`

	offset := (page - 1) * limit

	albumIDs := map[uuid.UUID]bool{}

	args1 := append(fq.Args, limit, offset)
	rows, err := d.ext.QueryContext(ctx, query1, args1...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var td types.TrackDetailedNew

		td.Mediafile = &types.MediaFile{}

		err := rows.Scan(
			&td.Track.ID,
			&td.Track.Title,
			&td.Track.MusicBrainzID,
			&td.Track.Genre,
			&td.Track.Comment,
			&td.Track.MediaFile,
			&td.Track.CreatedAt,
			&td.Track.UpdatedAt,

			&td.Album.ID,
			&td.Album.MusicBrainzID,
			&td.Album.Name,
			&td.Album.SortName,
			&td.Album.ReleaseDate,
			&td.Album.Tracks,
			&td.Album.Discs,
			&td.Album.Description,
			&td.Album.Owner,
			&td.Album.Public,
			&td.Album.CreatedAt,
			&td.Album.UpdatedAt,

			&td.AlbumTrack.ID,
			&td.AlbumTrack.AlbumID,
			&td.AlbumTrack.TrackID,
			&td.AlbumTrack.DiscNumber,
			&td.AlbumTrack.TrackNumber,
			&td.AlbumTrack.CreatedAt,
			&td.AlbumTrack.UpdatedAt,

			&td.Mediafile.ID,
			&td.Mediafile.DurationMs,
			&td.Mediafile.Bitrate,
			&td.Mediafile.SampleRate,
			&td.Mediafile.FileSize,
			&td.Mediafile.Codec,
			&td.Mediafile.Checksum,
			&td.Mediafile.CreatedAt,
			&td.Mediafile.UpdatedAt,

			&td.Analysis.ID,
			&td.Analysis.BPM,
			&td.Analysis.Key,
			&td.Analysis.KeyScale,
			&td.Analysis.KeyConfidence,
			&td.Analysis.Loudness,
			&td.Analysis.Energy,
			&td.Analysis.Danceability,
			&td.Analysis.MoodHappy,
			&td.Analysis.MoodSad,
			&td.Analysis.MoodAggresive,
			&td.Analysis.MoodCalm,
			&td.Analysis.CreatedAt,
			&td.Analysis.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}

		results = append(results, td)
		albumIDs[td.Album.ID] = true
	}

	ids := []uuid.UUID{}
	for id := range albumIDs {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, 0, 0, nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query2 := `
     SELECT aa.album_id, aa.role, ar.id, ar.name, ar.sort_name
     FROM album_artists aa
     JOIN artists ar ON ar.id = aa.artist_id
     WHERE aa.album_id IN (` + placeholders + `)`

	args := make([]any, len(ids))
	for i, v := range ids {
		args[i] = v
	}

	rows2, err := d.ext.QueryContext(ctx, query2, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows2.Close()

	artistMap := map[uuid.UUID][]types.Artist{}

	for rows2.Next() {
		var albumID uuid.UUID
		var a types.Artist

		err := rows2.Scan(&albumID, &a.Role, &a.ID, &a.Name, &a.SortName)
		if err != nil {
			return nil, 0, 0, err
		}

		artistMap[albumID] = append(artistMap[albumID], a)
	}

	for i := range results {
		results[i].Artists = artistMap[results[i].Album.ID]
	}

	countQuery := `
     SELECT COUNT(DISTINCT t.id)
     FROM track_analyses ta
     JOIN tracks t ON t.id = ta.id
     JOIN album_tracks at ON at.track_id = t.id
     JOIN albums al ON al.id = at.album_id
     LEFT JOIN media_files mf ON mf.id = t.media_file
`

	if len(fq.Joins) > 0 {
		countQuery += "\n" + strings.Join(fq.Joins, "\n")
	}

	if len(fq.Where) > 0 {
		countQuery += "\nWHERE " + strings.Join(fq.Where, " AND ")
	}

	err = sqlx.GetContext(ctx, d.ext, &totalItems, countQuery, args1...)
	if err != nil {
		return nil, 0, 0, err
	}

	if limit > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return results, totalPages, totalItems, nil
}
