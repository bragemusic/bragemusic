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
