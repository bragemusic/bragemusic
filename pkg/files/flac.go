package files

import (
	"io"
	"math"
	"os"

	"github.com/bragemusic/core/pkg/types"
	"github.com/mewkiz/flac"
)

type FileFlac struct {
	filesize int64
	stream   *flac.Stream
	buf      []byte
}

func (f FileFlac) SampleRate() int {
	return int(f.stream.Info.SampleRate)
}

func (f FileFlac) NChannels() int {
	return int(f.stream.Info.NChannels)
}

func (f FileFlac) Bitrate() int {
	return int(float32(f.filesize*8) / float32(f.DurationMS()/1000))
}

func (f FileFlac) FileSize() int64 {
	return f.filesize
}

func (f FileFlac) DurationMS() int64 {
	return (int64(f.stream.Info.NSamples) / int64(f.stream.Info.SampleRate)) * 1000
}

func (f *FileFlac) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		// Serve leftover buffer first
		if len(f.buf) > 0 {
			copied := copy(p[n:], f.buf)
			f.buf = f.buf[copied:]
			n += copied
			continue
		}

		// Decode next frame
		frame, err := f.stream.ParseNext()
		if err != nil {
			if err == io.EOF && n > 0 {
				return n, nil
			}
			return n, err
		}

		samples := len(frame.Subframes[0].Samples)
		channels := f.NChannels()
		tmp := make([]byte, samples*channels*4) // 4 bytes per float32

		// Compute max value for normalization
		bits := f.stream.Info.BitsPerSample
		maxVal := float32((int(1) << bits) - 1) // integer math, then convert to float32

		i := 0
		for s := 0; s < samples; s++ {
			for c := 0; c < channels; c++ {
				raw := frame.Subframes[c].Samples[s]
				val := float32(raw) / maxVal // normalize to [-1,1]

				bits := math.Float32bits(val)
				tmp[i] = byte(bits)
				tmp[i+1] = byte(bits >> 8)
				tmp[i+2] = byte(bits >> 16)
				tmp[i+3] = byte(bits >> 24)
				i += 4
			}
		}

		f.buf = tmp
	}
	return n, nil
}

func ParseFlac(f *os.File) (types.AudioFile, error) {
	ff := FileFlac{}

	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	flacFile, err := flac.New(f)
	if err != nil {
		return nil, err
	}

	ff.filesize = fstat.Size()
	ff.stream = flacFile

	// _, err = f.Seek(0, 0)
	// if err != nil {
	// 	return nil, err
	// }

	return &ff, nil
}
