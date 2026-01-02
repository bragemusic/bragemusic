package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Album struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	MusicBrainzID *string    `db:"musicbrainz_id" json:"musicbrainz_id"`
	Name          string     `db:"name" json:"name"`
	SortName      string     `db:"sort_name" json:"sort_name"`
	ReleaseDate   *time.Time `db:"release_date" json:"release_date,omitempty"`
	Tracks        *int       `db:"tracks" json:"tracks,omitempty"`
	Discs         *int       `db:"discs" json:"discs,omitempty"`
	Description   *string    `db:"description" json:"description,omitempty"`
	Owner         string     `db:"owner" json:"owner"`
	Public        *bool      `db:"public" json:"public,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

type AlbumEnhanced struct {
	ID            string     `db:"id" json:"id"`
	MusicBrainzID *string    `db:"musicbrainz_id" json:"musicbrainz_id"`
	Name          string     `db:"name" json:"name"`
	SortName      string     `db:"sort_name" json:"sort_name"`
	ArtistID      string     `db:"artist_id" json:"artist_id"`
	ArtistName    string     `db:"artist_name" json:"artist_name"`
	ReleaseDate   *time.Time `db:"release_date" json:"release_date,omitempty"`
	Tracks        *int       `db:"tracks" json:"tracks,omitempty"`
	Discs         *int       `db:"discs" json:"discs,omitempty"`
	Description   *string    `db:"description" json:"description,omitempty"`
	Owner         string     `db:"owner" json:"owner"`
	Public        *bool      `db:"public" json:"public,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}
