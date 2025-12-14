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

func (d Database) GetAuthIdentityForUser(ctx context.Context, userID uuid.UUID) (ai types.AuthIdentity, err error) {
	query := `
        SELECT *
        FROM auth_identities
        WHERE user_id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &ai, query, userID)
	if err != nil {
		return types.AuthIdentity{}, err
	}

	return
}

func (d Database) CreateAuthIdentity(ctx context.Context, ai types.AuthIdentity) (uuid.UUID, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	ai.ID = uid

	now := time.Now()
	ai.CreatedAt = now
	ai.UpdatedAt = now

	query := `
        INSERT INTO auth_identities(
            id, user_id, provider, provider_user_id, email, created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?);
    `

	_, err = d.ext.ExecContext(
		ctx,
		query,
		ai.ID,
		ai.UserID,
		ai.Provider,
		ai.ProviderUserID,
		ai.Email,
		ai.CreatedAt,
		ai.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return ai.ID, nil
}

func (d Database) UpdateAuthIdentity(ctx context.Context, ai types.AuthIdentity) error {
	ai.UpdatedAt = time.Now()
	query := `
        UPDATE auth_identities SET
            provider = :provider,
            provider_user_id = :provider_user_id,
            email = :email,
            updated_at = :updated_at
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, ai)
	return err
}
