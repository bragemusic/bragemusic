package database

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

type AuthFace interface {
	Begin(ctx context.Context) (DatabaseFace, error)
	driver.Tx

	UserExistsByID(ctx context.Context, ID uuid.UUID) (bool, error)
	CreateUser(ctx context.Context, user types.User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user types.User) error
	RemoveUser(ctx context.Context, userID uuid.UUID) error

	CreateAuthIdentity(ctx context.Context, ai types.AuthIdentity) (uuid.UUID, error)
	GetAuthIdentityForUser(ctx context.Context, userID uuid.UUID) (ai types.AuthIdentity, err error)
	UpdateAuthIdentity(ctx context.Context, ai types.AuthIdentity) error

	GetLocalCredentialsForUser(ctx context.Context, userID uuid.UUID) (lc types.LocalCredentials, err error)
	CreateLocalCredentials(ctx context.Context, lc types.LocalCredentials) error
	UpdateLocalCredentials(ctx context.Context, lc types.LocalCredentials) error

	CreateUserScope(ctx context.Context, us types.UserScope) error
	UserScopeExists(ctx context.Context, userID uuid.UUID, role types.UserRole) (bool, error)
	RemoveUserScope(ctx context.Context, userID uuid.UUID, role types.UserRole) error
	ListUserRoles(ctx context.Context, userID uuid.UUID) (roles []types.UserRole, err error)

	CreateToken(ctx context.Context, t types.Token) (uuid.UUID, error)
	GetTokenFromHash(ctx context.Context, hash string) (token types.Token, err error)

	GetUserDetails(ctx context.Context, userID uuid.UUID) (user types.UserDetails, err error)
}

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

func (d Database) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	query := `
        DELETE FROM users WHERE id = :id;
    `

	res, err := d.ext.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return ErrNoRowDeleted
	}

	return nil
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

func (d Database) GetLocalCredentialsForUser(ctx context.Context, userID uuid.UUID) (lc types.LocalCredentials, err error) {
	query := `
        SELECT *
        FROM local_credentials
        WHERE user_id = ?
        LIMIT 1;
    `
	err = sqlx.GetContext(ctx, d.ext, &lc, query, userID)
	if err != nil {
		return types.LocalCredentials{}, err
	}

	return
}

func (d Database) CreateLocalCredentials(ctx context.Context, lc types.LocalCredentials) error {
	now := time.Now()
	lc.CreatedAt = now
	lc.UpdatedAt = now

	query := `
        INSERT INTO local_credentials (
            user_id, password_hash, created_at, updated_at
        )
        VALUES (?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		lc.UserID,
		lc.PasswordHash,
		lc.CreatedAt,
		lc.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) UpdateLocalCredentials(ctx context.Context, lc types.LocalCredentials) error {
	lc.UpdatedAt = time.Now()
	query := `
        UPDATE local_credentials SET
            password_hash = :password_hash,
            updated_at = :updated_at
        WHERE user_id = :user_id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, lc)
	return err
}

func (d Database) CreateUserScope(ctx context.Context, us types.UserScope) error {
	now := time.Now()
	us.CreatedAt = now
	us.UpdatedAt = now

	query := `
        INSERT INTO user_scopes (
            user_id, role, created_at, updated_at
        )
        VALUES (?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		us.UserID,
		us.Role,
		us.CreatedAt,
		us.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) UserScopeExists(ctx context.Context, userID uuid.UUID, role types.UserRole) (bool, error) {
	const query = `
        SELECT COUNT(1)
        FROM user_scopes
        WHERE
          user_id = ?
          AND role = ?
        ;
    `
	var count int
	err := d.ext.QueryRowxContext(ctx, query, userID, role).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d Database) RemoveUserScope(ctx context.Context, userID uuid.UUID, role types.UserRole) error {
	query := `
        DELETE FROM user_scopes WHERE user_id = ? AND role = ?;
    `

	res, err := d.ext.ExecContext(ctx, query, userID, role)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return ErrNoRowDeleted
	}

	return nil
}

func (d Database) ListUserRoles(ctx context.Context, userID uuid.UUID) (roles []types.UserRole, err error) {
	query := `
        SELECT role
        FROM user_scopes
        WHERE user_id = ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &roles, query, userID)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) CreateToken(ctx context.Context, t types.Token) (uuid.UUID, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	t.ID = uid

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := `
        INSERT INTO tokens (
            id,
            user_id,
            type,
            name,
            hash,
            scopes,
            expires_at,
            last_used_at,
            created_at,
            updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err = d.ext.ExecContext(
		ctx,
		query,
		t.ID,
		t.UserID,
		t.Type,
		t.Name,
		t.Hash,
		t.Scopes,
		t.ExpiresAt,
		t.LastUsedAt,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return t.ID, nil
}

func (d Database) GetTokenFromHash(ctx context.Context, hash string) (token types.Token, err error) {
	query := `
        SELECT *
        FROM tokens
        WHERE hash = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &token, query, hash)
	if err != nil {
		return types.Token{}, err
	}

	return
}

func (d Database) GetUserDetails(ctx context.Context, userID uuid.UUID) (user types.UserDetails, err error) {
	// 1. Fetch user + primary provider
	queryUser := `
		SELECT
			u.id,
			u.email,
			u.username,
			u.created_at,
			ai.provider
		FROM users u
		LEFT JOIN auth_identities ai
			ON ai.user_id = u.id
		WHERE u.id = ?
		ORDER BY ai.created_at ASC
		LIMIT 1;
	`

	err = sqlx.GetContext(ctx, d.ext, &user, queryUser, userID)
	if err != nil {
		return user, err
	}

	// 2. Fetch roles (scopes)
	queryRoles := `
		SELECT role
		FROM user_scopes
		WHERE user_id = ?;
	`

	err = sqlx.SelectContext(ctx, d.ext, &user.Roles, queryRoles, userID)
	if err != nil {
		return user, err
	}

	return user, nil
}
