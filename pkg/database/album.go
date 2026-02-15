package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddAlbum(ctx context.Context, a types.Album, userID uuid.UUID) (uuid.UUID, error) {
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

	err = d.addEntityEvent(ctx, a.ID, types.EntityEventCreate, types.EntityAlbum, userID)
	if err != nil {
		return uuid.UUID{}, err
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
	query := `
		SELECT
			a.id,
			a.musicbrainz_id,
			a.name,
			a.sort_name,
			a.release_date,
			a.description,
			a.owner,
			a.public,
			a.created_at,
			a.updated_at,

			-- track count
			(
				SELECT COUNT(*)
				FROM album_tracks at
				WHERE at.album_id = a.id
			) AS tracks,

			-- disc count
			(
				SELECT COUNT(DISTINCT disc_number)
				FROM album_tracks at
				WHERE at.album_id = a.id
			) AS discs

		FROM albums a
		JOIN album_artists aa ON aa.album_id = a.id
		JOIN artists ar ON ar.id = aa.artist_id

		WHERE normalize(a.name) = normalize(?)
		  AND normalize(ar.name) = normalize(?)

		LIMIT 1;
	`

	err = sqlx.GetContext(ctx, d.ext, &album, query, albumName, artistName)
	if err != nil {
		return types.Album{}, err
	}

	return album, nil
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

func (d Database) GetAlbumFromID(ctx context.Context, id uuid.UUID) (album types.Album, err error) {
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

func (d Database) ListAlbums(ctx context.Context) (albums []types.Album, err error) {
	query := `
        SELECT *
        FROM albums
        ;
    `

	err = sqlx.SelectContext(ctx, d.ext, &albums, query)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (d Database) ListAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy SortBy, sortOrder SortOrder) (albums []types.AlbumDetailed, err error) {
	sortByStr := ""

	switch sortBy {
	case SortByName:
		sortByStr = "sort_name"
	case SortByDate:
		sortByStr = "release_date"
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT
			al.id,
			al.musicbrainz_id,
			al.name,
			al.sort_name,
			al.release_date,
			al.description,
			al.owner,
			al.public,
			al.created_at,
			al.updated_at
		FROM albums al
		JOIN album_artists aa ON aa.album_id = al.id
		WHERE aa.artist_id = ?
		ORDER BY %s %s;
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

func (d Database) UpdateAlbum(ctx context.Context, a types.Album, userID uuid.UUID) error {
	a.UpdatedAt = time.Now()
	query := `
        UPDATE albums
        SET
            musicbrainz_id = :musicbrainz_id,
            name = :name,
            sort_name = :sort_name,
            release_date = :release_date,
            tracks = :tracks,
            discs = :discs,
            description = :description,
            owner = :owner,
            public = :public,
            updated_at = :updated_at
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, a)
	if err != nil {
		return err
	}

	err = d.addEntityEvent(ctx, a.ID, types.EntityEventUpdate, types.EntityAlbum, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (album types.AlbumDetailed, err error) {
	// 1. Album core info
	albumQuery := `
		SELECT
			id,
			musicbrainz_id,
			name,
			sort_name,
			release_date,
			description,
			owner,
			public,
			created_at,
			updated_at
		FROM albums
		WHERE id = ?
		LIMIT 1;
	`
	if err := sqlx.GetContext(ctx, d.ext, &album, albumQuery, albumID); err != nil {
		return album, err
	}

	// 2. Album artists
	artistQuery := `
		SELECT a.id, a.name
		FROM album_artists aa
		JOIN artists a ON a.id = aa.artist_id
		WHERE aa.album_id = ?
		ORDER BY aa.position;
	`
	type artistRow struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	var artistRows []artistRow
	if err := sqlx.SelectContext(ctx, d.ext, &artistRows, artistQuery, albumID); err != nil {
		return album, err
	}

	for _, a := range artistRows {
		album.ArtistIDs = append(album.ArtistIDs, a.ID)
		album.ArtistNames = append(album.ArtistNames, a.Name)
	}

	return album, nil
}

func (d Database) CountAlbums(ctx context.Context) (int, error) {
	const query = `
        SELECT COUNT(1)
        FROM albums;
    `

	var count int
	err := d.ext.QueryRowxContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
