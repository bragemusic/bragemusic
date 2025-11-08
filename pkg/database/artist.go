package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddArtist(ctx context.Context, a types.Artist) (string, error) {
	if a.ID == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		a.ID = uid.String()
	}

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	query := `
        INSERT INTO artists (
            id, musicbrainz_id , name, sort_name, country, year_started, year_ended, description,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		a.ID,
		a.MusicBrainzID,
		a.Name,
		a.SortName,
		a.Country,
		a.YearStarted,
		a.YearEnded,
		a.Description,
		a.CreatedAt,
		a.UpdatedAt,
	)
	if err != nil {
		return "", err
	}

	return a.ID, nil
}

func (d Database) UpdateArtist(ctx context.Context, a types.Artist) error {
	query := `
        UPDATE artists
        SET
            musicbrainz_id = :musicbrainz_id,
            name = :name,
            sort_name = :sort_name,
            country = :country,
            year_started = :year_started,
            year_ended = :year_ended,
            description = :description
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, a)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) GetArtistFromID(ctx context.Context, id string) (artist types.Artist, err error) {
	query := `
        SELECT *
        FROM artists
        WHERE id = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &artist, query, id)
	if err != nil {
		return types.Artist{}, err
	}

	return
}

func (d Database) GetArtistFromMbID(ctx context.Context, mbID string) (artist types.Artist, err error) {
	query := `
        SELECT *
        FROM artists
        WHERE musicbrainz_id = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &artist, query, mbID)
	if err != nil {
		return types.Artist{}, err
	}

	return
}

func (d Database) GetArtistFromName(ctx context.Context, name string) (artist types.Artist, err error) {
	query := `
       SELECT *
       FROM artists
       WHERE normalize(name) = normalize(?)
       LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &artist, query, name)
	if err != nil {
		return types.Artist{}, err
	}

	return
}

func (d Database) ListArtists(ctx context.Context) (artists []types.Artist, err error) {
	query := `
        SELECT *
        FROM artists
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &artists, query)
	if err != nil {
		return nil, err
	}

	return
}
