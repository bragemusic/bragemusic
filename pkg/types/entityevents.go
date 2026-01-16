package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	EntityType      string
	EntityEventType string
)

const (
	EntityAlbumArtist EntityType = "album_artist"

	EntityEventDelete EntityEventType = "delete"
)

type EntityEvent struct {
	ID         uuid.UUID       `db:"id" json:"id"`
	Type       EntityEventType `db:"event_type" json:"event_type"`
	EntityType EntityType      `db:"entity_type" json:"entity_type"`
	EventTime  time.Time       `db:"event_time" json:"event_time"`
}
