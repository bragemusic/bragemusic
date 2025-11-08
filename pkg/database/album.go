package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddAlbum(ctx context.Context, a types.Album) (string, error) {
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

	const query = `
		INSERT INTO albums (
			id, musicbrainz_id, name, sort_name, artist_id, release_date, tracks, discs,
			description, owner, public, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		a.ID,
		a.MusicBrainzID,
		a.Name,
		a.SortName,
		a.ArtistID,
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
		return "", err
	}

	return a.ID, nil
}

func (d Database) GetAlbumFromArtistAndName(ctx context.Context, artistName, albumName string) (album types.Album, err error) {
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

func (d Database) ListAlbumsByArtist(ctx context.Context, artistID string) (albums []types.Album, err error) {
	query := `
        SELECT *
        FROM albums
        WHERE artist_id = ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &albums, query, artistID)
	if err != nil {
		return nil, err
	}

	return
}
