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
	EntityArtist               EntityType = "artist"
	EntityAlbum                EntityType = "album"
	EntityTrack                EntityType = "track"
	EntityAlbumArtist          EntityType = "album_artist"
	EntityAlbumTrack           EntityType = "album_track"
	EntityMediaFile            EntityType = "media_file"
	EntityPlaylist             EntityType = "playlist"
	EntitySmartPlaylistContent EntityType = "smart_playlist_content"
	EntitySmartPlaylistArtist  EntityType = "smart_playlist_artist"
	EntityPlaylistTrack        EntityType = "playlist_track"
	EntityPlayHistoryItem      EntityType = "play_history_item"
	EntitySearchItem           EntityType = "search_item"
	EntityImport               EntityType = "import"
	EntityRating               EntityType = "rating"
	EntityLike                 EntityType = "like"
	EntityEntityEvent          EntityType = "entity_event"
	EntityDevice               EntityType = "device"
	EntityDeviceToken          EntityType = "device_token"
	EntityToken                EntityType = "token"
	EntityTrackAnalysis        EntityType = "track_analysis"
	EntityTrackArtist          EntityType = "track_artist"

	EntityEventDelete EntityEventType = "delete"
	EntityEventCreate EntityEventType = "create"
	EntityEventUpdate EntityEventType = "update"
)

type EntityEvent struct {
	ID         uuid.UUID       `db:"id" json:"id" ts_type:"string"`
	ItemID     uuid.UUID       `db:"item_id" json:"item_id" ts_type:"string"`
	UserID     uuid.UUID       `db:"user_id" json:"user_id" ts_type:"string"`
	Type       EntityEventType `db:"event_type" json:"event_type"`
	EntityType EntityType      `db:"entity_type" json:"entity_type"`
	EventTime  time.Time       `db:"event_time" json:"event_time" ts_type:"string"`
}

type EntityEvents []EntityEvent

func (es EntityEvents) LaterDeleteExists(e EntityEvent) bool {
	if e.Type == EntityEventDelete {
		return false
	}

	for _, ee := range es {
		if ee.EventTime.Before(e.EventTime) {
			continue
		}

		if ee.ItemID != e.ItemID {
			continue
		}

		if ee.Type != EntityEventDelete {
			continue
		}

		return true

	}

	return false
}
