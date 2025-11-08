package files

import (
	"fmt"
	"os"

	"github.com/bragemusic/core/pkg/types"
	"github.com/dhowden/tag"
)

func ParseAudioFile(f *os.File, filetype tag.FileType) (types.AudioFile, error) {
	switch filetype {
	case tag.FLAC:
		return ParseFlac(f)
	default:
		return nil, fmt.Errorf("unsupported audio format '%s'", filetype)
	}
}
