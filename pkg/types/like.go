package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Like struct {
	ID        uuid.UUID `db:"id" json:"id"`
	TrackID   uuid.UUID `db:"track_id" json:"track_id"`
	Owner     uuid.UUID `db:"owner" json:"owner"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
