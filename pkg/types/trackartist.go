package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type TrackArtist struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	TrackID   uuid.UUID  `db:"track_id" json:"album_id"`
	ArtistID  uuid.UUID  `db:"artist_id" json:"artist_id"`
	Role      ArtistRole `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
