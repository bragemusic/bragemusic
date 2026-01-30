package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Playlist struct {
	ID          uuid.UUID `db:"id" json:"id" ts_type:"string"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	Owner       uuid.UUID `db:"owner" json:"owner"`
	Public      bool      `db:"public" json:"public"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type PlaylistTrack struct {
	ID           uuid.UUID `db:"id" json:"id" ts_type:"string"`
	PlaylistID   uuid.UUID `db:"playlist_id" json:"playlist_id"`
	AlbumTrackID uuid.UUID `db:"album_track_id" json:"album_track_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
