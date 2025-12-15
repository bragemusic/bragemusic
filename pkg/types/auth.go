package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	UserRole     string
	AuthProvider string
)

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleRead  UserRole = "read"
	UserRoleWrite UserRole = "write"

	AuthLocal AuthProvider = "local"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
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
