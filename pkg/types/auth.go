package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	UserRole     string
	AuthProvider string
	TokenType    string
)

const (
	UserRoleAdmin         UserRole = "admin"
	UserRoleRead          UserRole = "read"
	UserRoleWrite         UserRole = "write"
	UserRoleImporterRead  UserRole = "importer:read"
	UserRoleImporterWrite UserRole = "importer:write"
	UserRoleUsersRead     UserRole = "users:read"
	UserRoleUsersCreate   UserRole = "users:create"
	UserRoleUsersUpdate   UserRole = "users:update"
	UserRoleUsersDelete   UserRole = "users:delete"

	AuthLocal AuthProvider = "local"

	TokenFrontendLong  TokenType = "frontend_long"
	TokenFrontendShort TokenType = "frontend_short"
	TokenAPI           TokenType = "api"
)

var AllUserRoles = []UserRole{
	UserRoleAdmin,
	UserRoleRead,
	UserRoleWrite,
	UserRoleImporterWrite,
	UserRoleUsersRead,
	UserRoleUsersCreate,
	UserRoleUsersUpdate,
	UserRoleUsersDelete,
}

type UserDetails struct {
	ID        uuid.UUID    `db:"id" json:"id" ts_type:"string"`
	Email     string       `db:"email" json:"email"`
	Username  string       `db:"username" json:"username"`
	Provider  AuthProvider `db:"provider" json:"provider"`
	Roles     []UserRole   `db:"role" json:"role"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
}

func (u UserDetails) HasRole(r UserRole) bool {
	for _, rr := range u.Roles {
		if rr == r {
			return true
		}
	}
	return false
}

type User struct {
	ID        uuid.UUID `db:"id" json:"id" ts_type:"string"`
	Email     string    `db:"email" json:"email"`
	Username  string    `db:"username" json:"username"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AuthIdentity struct {
	ID             uuid.UUID    `db:"id" json:"id"`
	UserID         uuid.UUID    `db:"user_id" json:"user_id"`
	Provider       AuthProvider `db:"provider" json:"provider"`
	ProviderUserID uuid.UUID    `db:"provider_user_id" json:"provider_user_id"`
	Email          string       `db:"email" json:"email"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at" json:"updated_at"`
}

type LocalCredentials struct {
	UserID       uuid.UUID `db:"user_id" json:"user_id"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserScope struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Role      UserRole  `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Token struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	Type       TokenType  `db:"type" json:"type"`
	Name       *string    `db:"name" json:"name,omitempty"`
	Hash       string     `db:"hash" json:"-"`
	Scopes     string     `db:"scopes" json:"scopes"`
	ExpiresAt  *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	LastUsedAt *time.Time `db:"last_used_at" json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

type TokenLimited struct {
	ID         uuid.UUID  `db:"id" json:"id" ts_type:"string"`
	Type       TokenType  `db:"type" json:"type"`
	Name       *string    `db:"name" json:"name,omitempty"`
	Scopes     string     `db:"scopes" json:"scopes"`
	ExpiresAt  *time.Time `db:"expires_at" json:"expires_at,omitempty" ts_type:"string"`
	LastUsedAt *time.Time `db:"last_used_at" json:"last_used_at,omitempty" ts_type:"string"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at" ts_type:"string"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at" ts_type:"string"`
}
