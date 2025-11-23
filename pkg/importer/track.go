package importer

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
)

func (i Importer) updateTrackData(track types.Track, audioFile types.AudioFile, filename string, filetype tag.FileType) types.Track {
	track.DurationMS = utils.Ptr(audioFile.DurationMS())
	track.Bitrate = utils.Ptr(audioFile.Bitrate())
	track.SampleRate = utils.Ptr(audioFile.SampleRate())
	track.FilePath = filename
	track.FileSize = utils.Ptr(audioFile.FileSize())
	track.MimeType = utils.Ptr(string(filetype))

	return track
}
