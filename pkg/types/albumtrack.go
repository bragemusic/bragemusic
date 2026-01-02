package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type AlbumTrack struct {
	AlbumID     uuid.UUID `db:"album_id" json:"album_id"`
	TrackID     uuid.UUID `db:"track_id" json:"track_id"`
	DiscNumber  int       `db:"disc_number" json:"disc_number"`
	TrackNumber int       `db:"track_number" json:"track_number"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
