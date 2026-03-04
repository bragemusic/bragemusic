package database

import (
	"context"
	"strings"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddLike(ctx context.Context, l types.Like) (uuid.UUID, error) {
	if l.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		l.ID = uid
	}

	if l.CreatedAt.IsZero() {
		now := time.Now()
		l.CreatedAt = now
		l.UpdatedAt = now
	}

	query := `
        INSERT INTO likes (
            id, track_id, owner,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		l.ID,
		l.TrackID,
		l.Owner,
		l.CreatedAt,
		l.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = d.addEntityEvent(ctx, l.ID, types.EntityEventCreate, types.EntityLike, l.Owner)
	if err != nil {
		return uuid.UUID{}, err
	}

	return l.ID, nil
}

func (d Database) DeleteLike(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM likes
		WHERE
			id = ?;
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return err
	}

	err = d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityLike, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) GetLikeID(ctx context.Context, trackID, userID uuid.UUID) (uuid.UUID, error) {
	query := `
        SELECT id
        FROM likes
        WHERE track_id = ?
          AND owner = ?
        LIMIT 1;
    `

	var idStr string
	err := sqlx.GetContext(ctx, d.ext, &idStr, query, trackID, userID)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (d Database) HasLike(ctx context.Context, trackID, userID uuid.UUID) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM likes
        WHERE owner = ?
        AND track_id = ?
;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, userID, trackID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) ListLikedTracksDetailed(ctx context.Context, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
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
   	   FROM likes l
   	   JOIN tracks t ON t.id = l.track_id
   	   LEFT JOIN album_tracks at ON at.track_id = t.id
   	   LEFT JOIN albums al ON al.id = at.album_id
   	   LEFT JOIN (
   	   	SELECT track_id, COUNT(*) AS play_count
   	   	FROM play_history
   	   	GROUP BY track_id
   	   ) tp ON tp.track_id = t.id
   	   WHERE l.owner = ?
   	   ORDER BY l.created_at ASC;
	`

	if err := sqlx.SelectContext(ctx, d.ext, &tracks, tracksQuery, userID); err != nil {
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

func (d Database) attachTrackLike(ctx context.Context, tracks []types.TrackDetailed, userID uuid.UUID) error {
	if len(tracks) == 0 {
		return nil
	}

	args := make([]any, 0, len(tracks)+1)
	args = append(args, userID.String())

	placeholders := make([]string, len(tracks))
	for i, t := range tracks {
		placeholders[i] = "?"
		args = append(args, t.ID.String())
	}

	query := `
		SELECT track_id
		FROM likes
		WHERE owner = ?
		AND track_id IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := d.ext.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	liked := make(map[string]struct{}, len(tracks))

	for rows.Next() {
		var trackID string
		if err := rows.Scan(&trackID); err != nil {
			return err
		}
		liked[trackID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for i := range tracks {
		_, tracks[i].Liked = liked[tracks[i].ID.String()]
	}

	return nil
}
