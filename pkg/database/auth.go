package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) UserExistsByID(ctx context.Context, ID uuid.UUID) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM users
        WHERE id = ?;
    `
	var count int
	err := d.ext.QueryRowxContext(ctx, query, ID.String()).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) CreateUser(ctx context.Context, user types.User) (uuid.UUID, error) {
	if user.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		user.ID = uid
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
        INSERT INTO users (
            id, email, username, created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (d Database) UpdateUser(ctx context.Context, user types.User) error {
	user.UpdatedAt = time.Now()
	query := `
        UPDATE users SET
            email = :email,
            username = :username,
            updated_at = :updated_at
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, user)
	return err
}
