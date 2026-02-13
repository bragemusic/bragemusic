package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	ImportState string
	ImportType  string
)

const (
	ImportStateNotStarted ImportState = "not_started"
	ImportStateRunning    ImportState = "running"
	ImportStateFinished   ImportState = "finished"
	ImportStateError      ImportState = "error"

	ImportTypeAlbum ImportType = "album"
	ImportTypeTrack ImportType = "track"
)

type ImportAlbum struct {
	MusicbrainzID *string `json:"musicbrainz_id"`
}

type Import struct {
	ID            uuid.UUID   `db:"id" json:"id"`
	MusicBrainzID *string     `db:"musicbrainz_id" json:"musicbrainz_id"`
	Owner         uuid.UUID   `db:"owner" json:"owner"`
	Filename      string      `db:"filename" json:"filename"`
	Type          ImportType  `db:"type" json:"type"`
	State         ImportState `db:"state" json:"state"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
}
