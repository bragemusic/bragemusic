package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Artist struct {
	ID            uuid.UUID `db:"id" json:"id"`
	MusicBrainzID *string   `db:"musicbrainz_id" json:"musicbrainz_id"`
	Name          string    `db:"name" json:"name"`
	SortName      string    `db:"sort_name" json:"sort_name"`
	Country       *string   `db:"country" json:"country,omitempty"`
	YearStarted   *int      `db:"year_started" json:"year_started,omitempty"`
	YearEnded     *int      `db:"year_ended" json:"year_ended,omitempty"`
	Description   *string   `db:"description" json:"description,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type ArtistDetailed struct {
	Artist
	AlbumCount int `json:"album_count"`
	TrackCount int `json:"track_count"`
}
