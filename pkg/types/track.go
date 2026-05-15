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
	ID            uuid.UUID     `db:"id" json:"id" ts_type:"string"` // UUID
	Title         string        `db:"title" json:"title"`
	AlbumID       string        `db:"album_id" json:"album_id,omitempty"`
	AlbumName     string        `db:"album_name" json:"album_name,omitempty"`
	ArtistIDs     []string      `db:"artist_ids" json:"artist_ids,omitempty"`
	ArtistNames   []string      `db:"artist_names" json:"artist_names,omitempty"`
	MusicBrainzID *string       `db:"musicbrainz_id" json:"musicbrainz_id"`
	TrackNumber   int           `db:"track_number" json:"track_number,omitempty"`
	DiscNumber    int           `db:"disc_number" json:"disc_number,omitempty"`
	Genre         *string       `db:"genre" json:"genre,omitempty"`
	Comment       *string       `db:"comment" json:"comment,omitempty"`
	MediaFile     *MediaFile    `db:"media_file" json:"media_file"`
	PlayCount     int           `db:"play_count" json:"play_count"`
	ContextID     *uuid.UUID    `db:"context_id" json:"context_id" ts_type:"string"`
	Rating        *float64      `json:"rating"`
	UserRating    *int          `json:"user_rating"`
	Liked         bool          `json:"liked"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
	Analysis      TrackAnalysis `json:"analysis"`
}

type TrackUpdate struct {
	Track
	DiscNumber  int       `json:"disc_number"`
	TrackNumber int       `json:"track_number"`
	AlbumID     uuid.UUID `json:"album_id"`
	// Artists     []uuid.UUID `json:"artists"`
}

type TrackDetailedNew struct {
	Track      Track         `json:"track"`
	Album      Album         `json:"album"`
	Artists    []Artist      `json:"artists"`
	AlbumTrack AlbumTrack    `json:"album_track"`
	Mediafile  *MediaFile    `json:"media_file"`
	Analysis   TrackAnalysis `json:"analysis"`
}

type TrackFilter struct {
	BPM     *FilterUpperLower[int] `json:"bpm"`
	Mood    FilterMood             `json:"mood"`
	Artists *[]uuid.UUID           `json:"artists" ts_type:"string[]"`
}

type FilterMood struct {
	Aggressive *float64 `json:"aggressive"`
	Calm       *float64 `json:"calm"`
	Happy      *float64 `json:"happy"`
	Sad        *float64 `json:"sad"`
}
type FilterUpperLower[T any] struct {
	Upper T `json:"upper"`
	Lower T `json:"lower"`
}
