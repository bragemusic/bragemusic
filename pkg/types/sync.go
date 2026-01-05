package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	SyncItemState string
	SyncItemType  string
)

const (
	SiStateNotStarted SyncItemState = "NotStarted"
	SiStateRunning    SyncItemState = "Running"
	SiStateFinished   SyncItemState = "Finished"

	SiTypeMediaFile SyncItemType = "MediaFile"
)

type SyncState struct {
	Time    time.Time `json:"time"`
	Artists []string  `json:"artists"`
	Albums  []string  `json:"albums"`
	Tracks  []string  `json:"tracks"`
}

type DBSyncState struct {
	ID             string    `db:"id" json:"id"`
	SyncedAt       time.Time `db:"synced_at" json:"synced_at"`
	ArtistsCreated int       `db:"artists_created" json:"artists_created"`
	ArtistsUpdated int       `db:"artists_updated" json:"artists_updated"`
	AlbumsCreated  int       `db:"albums_created" json:"albums_created"`
	AlbumsUpdated  int       `db:"albums_updated" json:"albums_updated"`
	TracksCreated  int       `db:"tracks_created" json:"tracks_created"`
	TracksUpdated  int       `db:"tracks_updated" json:"tracks_updated"`
}

type SyncItem struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	SyncID    uuid.UUID     `db:"sync_id" json:"sync_id"`
	ItemID    uuid.UUID     `db:"item_id" json:"item_id"`
	Type      SyncItemType  `db:"type" json:"type"`
	State     SyncItemState `db:"state" json:"state"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at"`
}

type PlayHistorySyncState struct {
	Time        time.Time     `json:"time"`
	RemoteItems []PlayHistory `json:"remote_items"`
}
