package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	EntityType      string
	EntityEventType string
)

func (e EntityType) P() *EntityType {
	return &e
}

const (
	EntityArtist          EntityType = "artist"
	EntityAlbum           EntityType = "album"
	EntityTrack           EntityType = "track"
	EntityAlbumArtist     EntityType = "album_artist"
	EntityAlbumTrack      EntityType = "album_track"
	EntityMediaFile       EntityType = "media_file"
	EntityPlaylist        EntityType = "playlist"
	EntityPlaylistTrack   EntityType = "playlist_track"
	EntityPlayHistoryItem EntityType = "play_history_item"
	EntitySearchItem      EntityType = "search_item"
	EntityImport          EntityType = "import"

	EntityEventDelete EntityEventType = "delete"
)

type EntityEvent struct {
	ID         uuid.UUID       `db:"id" json:"id"`
	Type       EntityEventType `db:"event_type" json:"event_type"`
	EntityType EntityType      `db:"entity_type" json:"entity_type"`
	EventTime  time.Time       `db:"event_time" json:"event_time"`
}
