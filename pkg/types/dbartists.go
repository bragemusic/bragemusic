package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type ArtistRole string

const (
	ArPrimary         ArtistRole = "primary"
	ArFeatured        ArtistRole = "featured"
	ArComposer        ArtistRole = "composer"
	ArConductor       ArtistRole = "conductor"
	ArProducer        ArtistRole = "producer"
	ArRemixer         ArtistRole = "remixer"
	ArExecutive       ArtistRole = "executive"
	ArPerformer       ArtistRole = "performer"
	ArVocalist        ArtistRole = "vocalist"
	ArLyricist        ArtistRole = "lyricist"
	ArInstrumentalist ArtistRole = "instrumentalist"
)

type ArtistBase struct {
	MusicBrainzID *string `db:"musicbrainz_id" json:"musicbrainz_id"`
	Name          string  `db:"name" json:"name" required:"true"`
	SortName      string  `db:"sort_name" json:"sort_name" required:"true"`
	Country       *string `db:"country" json:"country,omitempty"`
	YearStarted   *int    `db:"year_started" json:"year_started,omitempty"`
	YearEnded     *int    `db:"year_ended" json:"year_ended,omitempty"`
	Description   *string `db:"description" json:"description,omitempty"`
}

type Artist struct {
	ID uuid.UUID `db:"id" json:"id" ts_type:"string"`
	ArtistBase
	Role      ArtistRole `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type ArtistMinimal struct {
	ID       uuid.UUID  `db:"id" json:"id" ts_type:"string"`
	Name     string     `db:"name" json:"name" required:"true"`
	SortName string     `db:"sort_name" json:"sort_name" required:"true"`
	Role     ArtistRole `db:"role" json:"role"`
}

type ArtistDetailed struct {
	Artist
	AlbumCount int `json:"album_count"`
	TrackCount int `json:"track_count"`
}
