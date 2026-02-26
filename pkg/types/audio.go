package types

import "io"

type MediaStream interface {
	io.ReadSeekCloser
	Size() (int64, error)
}

type AudioFile interface {
	io.Reader
	SampleRate() int
	NChannels() int
	Bitrate() int
	FileSize() int64
	DurationMS() int64
}
