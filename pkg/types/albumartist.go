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

type AlbumArtist struct {
	AlbumID   uuid.UUID  `db:"album_id" json:"album_id"`
	ArtistID  uuid.UUID  `db:"artist_id" json:"artist_id"`
	Role      ArtistRole `db:"role" json:"role"`
	Position  int        `db:"position" json:"position"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
