package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Track struct {
	ID            uuid.UUID  `db:"id" json:"id"` // UUID
	Title         string     `db:"title" json:"title"`
	MusicBrainzID *string    `db:"musicbrainz_id" json:"musicbrainz_id"`
	Genre         *string    `db:"genre" json:"genre,omitempty"`
	Comment       *string    `db:"comment" json:"comment,omitempty"`
	MediaFile     *uuid.UUID `db:"media_file" json:"media_file"`
	TitleMatch    float32    `db:"-" json:"-"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

type TrackDetailed struct {
	ID            string     `db:"id" json:"id"` // UUID
	Title         string     `db:"title" json:"title"`
	AlbumID       string     `db:"album_id" json:"album_id,omitempty"`
	AlbumName     string     `db:"album_name" json:"album_name,omitempty"`
	ArtistIDs     []string   `db:"artist_ids" json:"artist_ids,omitempty"`
	ArtistNames   []string   `db:"artist_names" json:"artist_names,omitempty"`
	MusicBrainzID *string    `db:"musicbrainz_id" json:"musicbrainz_id"`
	TrackNumber   int        `db:"track_number" json:"track_number,omitempty"`
	DiscNumber    int        `db:"disc_number" json:"disc_number,omitempty"`
	Genre         *string    `db:"genre" json:"genre,omitempty"`
	Comment       *string    `db:"comment" json:"comment,omitempty"`
	MediaFile     *MediaFile `db:"media_file" json:"media_file"`
	PlayCount     int        `db:"play_count" json:"play_count"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

type TrackUpdate struct {
	Track
	DiscNumber  int `json:"disc_number"`
	TrackNumber int `json:"track_number"`
	// Artists     []uuid.UUID `json:"artists"`
}
