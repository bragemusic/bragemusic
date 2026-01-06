package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddMediaFile(ctx context.Context, mf types.MediaFile) (uuid.UUID, error) {
	if mf.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		mf.ID = uid
	}

	now := time.Now()
	mf.CreatedAt = now
	mf.UpdatedAt = now

	query := `
        INSERT INTO media_files (
            id, duration_ms, bitrate, sample_rate,
            file_size, codec, checksum,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		mf.ID,
		mf.DurationMs,
		mf.Bitrate,
		mf.SampleRate,
		mf.FileSize,
		mf.Codec,
		mf.Checksum,
		mf.CreatedAt,
		mf.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return mf.ID, nil
}

func (d Database) GetMediaFileFromChecksum(ctx context.Context, cs string) (mf types.MediaFile, err error) {
	query := `
        SELECT *
        FROM media_files
        WHERE checksum = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &mf, query, cs)
	if err != nil {
		return types.MediaFile{}, err
	}

	return
}

func (d Database) attachMediaFiles(ctx context.Context, tracks []types.TrackDetailed) error {
	trackIDs := make([]string, 0, len(tracks))
	trackIndex := make(map[string]*types.TrackDetailed)

	for i := range tracks {
		trackIDs = append(trackIDs, tracks[i].ID)
		trackIndex[tracks[i].ID] = &tracks[i]
	}

	var rows []struct {
		TrackID string `db:"track_id"`
		types.MediaFile
	}

	query := `
		SELECT
			t.id AS track_id,
			mf.id,
			mf.duration_ms,
			mf.bitrate,
			mf.sample_rate,
			mf.file_size,
			mf.codec,
			mf.checksum,
			mf.created_at,
			mf.updated_at
		FROM tracks t
		JOIN media_files mf ON mf.id = t.media_file
		WHERE t.id IN (?);
	`

	query, args, err := sqlx.In(query, trackIDs)
	if err != nil {
		return err
	}
	query = d.ext.Rebind(query)

	if err := sqlx.SelectContext(ctx, d.ext, &rows, query, args...); err != nil {
		return err
	}

	for _, r := range rows {
		if t := trackIndex[r.TrackID]; t != nil {
			mf := r.MediaFile
			t.MediaFile = &mf
		}
	}

	return nil
}

func (d Database) ListUpdatedMediaFiles(ctx context.Context, since time.Time) (mediaFileIDs []uuid.UUID, err error) {
	query := `
        SELECT id
        FROM media_files
        WHERE
          created_at > ?
          OR
          updated_at > ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &mediaFileIDs, query, since, since)
	if err != nil {
		return nil, err
	}

	return
}
