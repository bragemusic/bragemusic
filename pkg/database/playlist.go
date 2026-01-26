package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddPlaylist(ctx context.Context, p types.Playlist) (uuid.UUID, error) {
	if p.Owner == uuid.Nil {
		return uuid.Nil, ErrNoUser
	}

	if p.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		p.ID = uid
	}

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	const query = `
		INSERT INTO playlists (
			id, name, description, owner, public,
            created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		p.ID,
		p.Name,
		p.Description,
		p.Owner,
		p.Public,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (d Database) DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		DELETE FROM playlists
		WHERE
			id = ?
          AND
            owner = ?;
	`

	_, err := d.ext.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	return d.addEntityEvent(ctx, id, types.EntityEventDelete, types.EntityPlaylist)
}

func (d Database) GetPlaylist(ctx context.Context, ID, userID uuid.UUID) (plist types.Playlist, err error) {
	query := `
        SELECT *
        FROM playlists
        WHERE id = ?
          AND (owner = ? OR public = 1)
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &plist, query, ID, userID)
	if err != nil {
		return types.Playlist{}, err
	}

	return
}

func (d Database) ListPlaylists(ctx context.Context, userID uuid.UUID, includePublic bool, sortBy SortBy, sortOrder SortOrder) (playlists []types.Playlist, err error) {
	sortByStr := ""

	switch sortBy {
	case SortByDate:
		sortByStr = "created_at"
	case SortByName:
		sortByStr = "name"
	case SortByPlayCount:
		sortByStr = "name"
	}

	publicStr := ""
	if includePublic {
		publicStr = "OR public = 1"
	}

	query := fmt.Sprintf(`
        SELECT *
        FROM playlists
        WHERE
          owner = ?
          %s
        ORDER BY %s %s
        ;
    `, publicStr, sortByStr, sortOrder)
	err = sqlx.SelectContext(ctx, d.ext, &playlists, query, userID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListUpdatedPlaylists(ctx context.Context, since time.Time, userID uuid.UUID) (plists []uuid.UUID, err error) {
	query := `
		SELECT
			id
		FROM playlists
		WHERE
            owner = ?
          AND
			(created_at > ?
			OR updated_at > ?);
    `
	err = sqlx.SelectContext(ctx, d.ext, &plists, query, userID, since, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) PlaylistExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
        SELECT COUNT(1)
        FROM playlists
        WHERE id = ?;
	`

	var exists bool
	err := sqlx.GetContext(ctx, d.ext, &exists, query, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (d Database) UpdatePlaylist(ctx context.Context, plist types.Playlist) error {
	plist.UpdatedAt = time.Now()
	query := `
        UPDATE playlists SET
            name = :name,
            description = :description,
            updated_at = :updated_at
        WHERE
            id = :id
            AND owner = :owner;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, plist)
	return err
}
