package types

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Codec string

const (
	CodecFlac Codec = "flac"
)

type MediaFile struct {
	ID         uuid.UUID `db:"id" json:"id"`
	DurationMs int64     `db:"duration_ms" json:"duration_ms"`
	Bitrate    int       `db:"bitrate" json:"bitrate"`
	SampleRate int       `db:"sample_rate" json:"sample_rate"`
	FileSize   int64     `db:"file_size" json:"file_size"`
	Codec      Codec     `db:"codec" json:"codec"`
	Checksum   string    `db:"checksum" json:"checksum"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (m MediaFile) Filename() string {
	return fmt.Sprintf("%s.%s", m.ID.String(), m.Codec)
}
