package types

import "time"

type Codec string

const (
	CodecFlac Codec = "flac"
)

type MediaFile struct {
	ID         string    `db:"id" json:"id"`
	DurationMs int64     `db:"duration_ms" json:"duration_ms"`
	Bitrate    int       `db:"bitrate" json:"bitrate"`
	SampleRate int       `db:"sample_rate" json:"sample_rate"`
	FilePath   string    `db:"file_path" json:"file_path"`
	FileSize   int64     `db:"file_size" json:"file_size"`
	Codec      Codec     `db:"codec" json:"codec"`
	Checksum   string    `db:"checksum" json:"checksum"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
