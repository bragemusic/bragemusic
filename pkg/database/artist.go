package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddArtist(ctx context.Context, a types.Artist) (uuid.UUID, error) {
	if a.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		a.ID = uid
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
		return uuid.Nil, err
	}

	return a.ID, nil
}

func (d Database) ArtistExists(ctx context.Context, ID string) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM artists
        WHERE id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, ID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) UpdateArtist(ctx context.Context, a types.Artist) error {
	a.UpdatedAt = time.Now()
	query := `
        UPDATE artists
        SET
            musicbrainz_id = :musicbrainz_id,
            name = :name,
            sort_name = :sort_name,
            country = :country,
            year_started = :year_started,
            year_ended = :year_ended,
            description = :description,
            updated_at = :updated_at
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

func (d Database) ListArtists(ctx context.Context, sortBy SortBy, sortOrder SortOrder) (artists []types.Artist, err error) {
	sortByStr := ""

	switch sortBy {
	case SortByDate:
		sortByStr = "created_at"
	case SortByName:
		sortByStr = "sort_name"
	case SortByPlayCount:
		sortByStr = "sort_name"
	}

	query := fmt.Sprintf(`
        SELECT *
        FROM artists
        ORDER BY %s %s
        ;
    `, sortByStr, sortOrder)
	err = sqlx.SelectContext(ctx, d.ext, &artists, query)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListUpdatedArtists(ctx context.Context, since time.Time) (artistIDs []string, err error) {
	query := `
        SELECT id
        FROM artists
        WHERE
          created_at > ?
          OR
          updated_at > ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &artistIDs, query, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListArtistsWithoutMeta(ctx context.Context) (artists []types.Artist, err error) {
	query := `
        SELECT *
        FROM artists
        WHERE
          musicbrainz_id IS NOT null
        AND
          description IS null
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &artists, query)
	if err != nil {
		return nil, err
	}

	return
}
