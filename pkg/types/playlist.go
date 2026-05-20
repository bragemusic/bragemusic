package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type PlaylistType string

const (
	PlaylistTypeStandard PlaylistType = "standard"
	PlaylistTypeSmart    PlaylistType = "smart"
)

type PlaylistBase struct {
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description,omitempty"`
	Public      bool    `db:"public" json:"public"`
}

type Playlist struct {
	ID uuid.UUID `db:"id" json:"id" ts_type:"string"`
	PlaylistBase
	Type      PlaylistType `db:"type" json:"type"`
	Owner     uuid.UUID    `db:"owner" json:"owner"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
}

type SmartPlaylist struct {
	Playlist
	Content SmartPlaylistContent
}

type PlaylistTrack struct {
	ID           uuid.UUID `db:"id" json:"id" ts_type:"string"`
	PlaylistID   uuid.UUID `db:"playlist_id" json:"playlist_id"`
	AlbumTrackID uuid.UUID `db:"album_track_id" json:"album_track_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type SmartPlaylistContent struct {
	ID             uuid.UUID    `db:"id" json:"id" ts_type:"string"`
	PlaylistID     uuid.UUID    `db:"playlist_id" json:"playlist_id"`
	BpmUpper       *int         `db:"bpm_upper" json:"bpm_upper"`
	BpmLower       *int         `db:"bpm_lower" json:"bpm_lower"`
	MoodAggressive *float64     `db:"mood_aggressive" json:"mood_aggressive"`
	MoodCalm       *float64     `db:"mood_calm" json:"mood_calm"`
	MoodHappy      *float64     `db:"mood_happy" json:"mood_happy"`
	MoodSad        *float64     `db:"mood_sad" json:"mood_sad"`
	Artists        *[]uuid.UUID `json:"artists" ts_type:"string[]"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (s SmartPlaylistContent) TrackFilter() TrackFilter {
	tf := TrackFilter{
		Mood: FilterMood{
			Aggressive: s.MoodAggressive,
			Calm:       s.MoodCalm,
			Happy:      s.MoodHappy,
			Sad:        s.MoodSad,
		},
		Artists: s.Artists,
	}

	if s.BpmLower != nil && s.BpmUpper != nil {
		tf.BPM = &FilterUpperLower[int]{
			Upper: *s.BpmUpper,
			Lower: *s.BpmLower,
		}
	}

	return tf
}

type SmartPlaylistArtist struct {
	ID        uuid.UUID `db:"id" json:"id" ts_type:"string"`
	ParentID  uuid.UUID `db:"parent_id" json:"parent_id"`
	ArtistID  uuid.UUID `db:"artist_id" json:"artist_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
