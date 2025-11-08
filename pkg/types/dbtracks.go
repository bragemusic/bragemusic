package types

import (
	"time"
)

type Track struct {
	ID            string    `db:"id" json:"id"` // UUID
	Title         string    `db:"title" json:"title"`
	AlbumID       *string   `db:"album_id" json:"album,omitempty"`
	MusicBrainzID *string   `db:"musicbrainz_id" json:"musicbrainz_id"`
	TrackArtist   *string   `db:"track_artist" json:"album_artist,omitempty"`
	TrackNumber   *int      `db:"track_number" json:"track_number,omitempty"`
	DiscNumber    *int      `db:"disc_number" json:"disc_number,omitempty"`
	Genre         *string   `db:"genre" json:"genre,omitempty"`
	Year          *int      `db:"year" json:"year,omitempty"`
	Composer      *string   `db:"composer" json:"composer,omitempty"`
	Comment       *string   `db:"comment" json:"comment,omitempty"`
	DurationMS    *int64    `db:"duration_ms" json:"duration_ms,omitempty"` // milliseconds
	Bitrate       *int      `db:"bitrate" json:"bitrate,omitempty"`         // kbps
	SampleRate    *int      `db:"sample_rate" json:"sample_rate,omitempty"` // Hz
	FilePath      string    `db:"file_path" json:"file_path"`
	FileSize      *int64    `db:"file_size" json:"file_size,omitempty"` // bytes
	MimeType      *string   `db:"mime_type" json:"mime_type,omitempty"`
	TitleMatch    float32   `db:"-" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type TrackEnhanced struct {
	ID            string    `db:"id" json:"id"` // UUID
	Title         string    `db:"title" json:"title"`
	AlbumID       *string   `db:"album_id" json:"album_id,omitempty"`
	AlbumName     *string   `db:"album_name" json:"album_name,omitempty"`
	ArtistID      *string   `db:"artist_id" json:"artist_id,omitempty"`
	ArtistName    *string   `db:"artist_name" json:"artist_name,omitempty"`
	MusicBrainzID *string   `db:"musicbrainz_id" json:"musicbrainz_id"`
	TrackArtist   *string   `db:"track_artist" json:"album_artist,omitempty"`
	TrackNumber   *int      `db:"track_number" json:"track_number,omitempty"`
	DiscNumber    *int      `db:"disc_number" json:"disc_number,omitempty"`
	Genre         *string   `db:"genre" json:"genre,omitempty"`
	Year          *int      `db:"year" json:"year,omitempty"`
	Composer      *string   `db:"composer" json:"composer,omitempty"`
	Comment       *string   `db:"comment" json:"comment,omitempty"`
	DurationMS    *int64    `db:"duration_ms" json:"duration_ms,omitempty"` // milliseconds
	Bitrate       *int      `db:"bitrate" json:"bitrate,omitempty"`         // kbps
	SampleRate    *int      `db:"sample_rate" json:"sample_rate,omitempty"` // Hz
	FilePath      string    `db:"file_path" json:"file_path"`
	FileSize      *int64    `db:"file_size" json:"file_size,omitempty"` // bytes
	MimeType      *string   `db:"mime_type" json:"mime_type,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
