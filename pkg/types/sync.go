package types

import "time"

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

type PlayHistorySyncState struct {
	Time        time.Time     `json:"time"`
	RemoteItems []PlayHistory `json:"remote_items"`
}
