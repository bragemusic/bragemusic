package types

import "io"

type AudioFile interface {
	io.Reader
	SampleRate() int
	NChannels() int
	Bitrate() int
	FileSize() int64
	DurationMS() int64
}
