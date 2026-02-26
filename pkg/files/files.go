package files

import (
	"fmt"

	"github.com/bragemusic/core/pkg/types"
)

func ParseAudioFile(f types.MediaStream, codec types.Codec) (types.AudioFile, error) {
	switch codec {
	case types.CodecFlac:
		return ParseFlac(f)
	default:
		return nil, fmt.Errorf("unsupported audio format '%s'", codec)
	}
}
