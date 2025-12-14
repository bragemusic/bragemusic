package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	UserScope    string
	AuthProvider string
)

const (
	UserScopeAdmin UserScope = "admin"
	UserScopeRead  UserScope = "read"
	UserScopeWrite UserScope = "write"

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

// CREATE TABLE IF NOT EXISTS auth_identities (
//     id                TEXT PRIMARY KEY,           -- UUID
//     user_id           TEXT NOT NULL,
//     provider          TEXT NOT NULL,        -- "local", "github", "authelia"
//     provider_user_id  TEXT NOT NULL,
//     email             TEXT,
//     created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
//     updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
