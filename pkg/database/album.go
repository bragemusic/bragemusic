package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddAlbum(ctx context.Context, a types.Album) (uuid.UUID, error) {
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

	const query = `
		INSERT INTO albums (
			id, musicbrainz_id, name, sort_name, release_date, tracks, discs,
			description, owner, public, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		a.ID,
		a.MusicBrainzID,
		a.Name,
		a.SortName,
		a.ReleaseDate,
		a.Tracks,
		a.Discs,
		a.Description,
		a.Owner,
		a.Public,
		a.CreatedAt,
		a.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return a.ID, nil
}

func (d Database) AlbumExists(ctx context.Context, ID string) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM albums
        WHERE id = ?;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query, ID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) GetAlbumFromArtistAndName(ctx context.Context, artistName, albumName string) (album types.Album, err error) {
	return types.Album{}, sql.ErrNoRows
	query := `
       SELECT a.*
       FROM albums a
       JOIN artists ar ON a.artist_id = ar.id
       WHERE normalize(a.name) = normalize(?)
         AND normalize(ar.name) = normalize(?)
       LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &album, query, albumName, artistName)
	if err != nil {
		return types.Album{}, err
	}

	return
}

func (d Database) GetAlbumFromMbID(ctx context.Context, mbID string) (album types.Album, err error) {
	query := `
        SELECT *
        FROM albums
        WHERE musicbrainz_id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &album, query, mbID)
	if err != nil {
		return types.Album{}, err
	}

	return
}

func (d Database) GetAlbumFromID(ctx context.Context, id string) (album types.Album, err error) {
	query := `
        SELECT *
        FROM albums
        WHERE id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &album, query, id)
	if err != nil {
		return types.Album{}, err
	}

	return
}

func (d Database) GetEnhancedAlbumFromID(ctx context.Context, id string) (album types.AlbumEnhanced, err error) {
	query := `
        SELECT
            al.*,
            ar.name  AS artist_name
        FROM albums al
        JOIN artists ar ON al.artist_id = ar.id
        WHERE al.id = ?;
    `
	err = sqlx.GetContext(ctx, d.ext, &album, query, id)
	if err != nil {
		return types.AlbumEnhanced{}, err
	}

	return
}

func (d Database) GetAlbumsByMbIDs(ctx context.Context, albumMbIds []string) ([]types.Album, error) {
	query, args, err := sqlx.In(`
        SELECT *
        FROM albums
        WHERE musicbrainz_id IN (?)
    `, albumMbIds)
	if err != nil {
		return nil, err
	}

	query = d.ext.Rebind(query)

	var albums []types.Album
	err = sqlx.SelectContext(ctx, d.ext, &albums, query, args...)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (d Database) ListAlbumsByArtist(ctx context.Context, artistID string, sortBy SortBy, sortOrder SortOrder) (albums []types.Album, err error) {
	sortByStr := ""

	switch sortBy {
	case SortByName:
		sortByStr = "sort_name"
	case SortByDate:
		sortByStr = "release_date"
	}

	query := fmt.Sprintf(`
        SELECT *
        FROM albums
        WHERE artist_id = ?
        ORDER BY %s %s
        ;
    `, sortByStr, sortOrder)
	err = sqlx.SelectContext(ctx, d.ext, &albums, query, artistID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListUpdatedAlbums(ctx context.Context, since time.Time) (albumIDs []string, err error) {
	query := `
        SELECT id
        FROM albums
        WHERE
          created_at > ?
          OR
          updated_at > ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &albumIDs, query, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) UpdateAlbum(ctx context.Context, a types.Album) error {
	query := `
        UPDATE albums
        SET
            musicbrainz_id = :musicbrainz_id,
            name = :name,
            sort_name = :sort_name,
            artist_id = :artist_id,
            release_date = :release_date,
            tracks = :tracks,
            discs = :discs,
            description = :description,
            owner = :owner,
            public = :public
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, a)
	if err != nil {
		return err
	}

	return nil
}
