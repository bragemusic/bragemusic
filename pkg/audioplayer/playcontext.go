package audioplayer

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type (
	PlayContextType string
	RepeatType      string
)

const (
	PlayContextAlbum     PlayContextType = "album"
	PlayContextTopTracks PlayContextType = "top_tracks"
	PlayContextPlaylist  PlayContextType = "playlist"

	RepeatOne RepeatType = "one"
	RepeatAll RepeatType = "all"
	RepeatOff RepeatType = "off"
)

type PlayContext struct {
	Type            PlayContextType       `json:"type"`
	RefID           uuid.UUID             `json:"ref_id" ts_type:"string"`
	Tracks          []types.TrackDetailed `json:"tracks"`
	Queue           []types.TrackDetailed `json:"queue"`
	CurrentTrackIdx int                   `json:"current_track_idx"`
	CurrentTrack    *types.TrackDetailed  `json:"current_track"`
	Shuffle         bool                  `json:"shuffle"`
	Repeat          RepeatType            `json:"repeat"`
	trackOrder      []int                 `json:"-"`
}

func (pc *PlayContext) PullFromQueue() *types.TrackDetailed {
	if len(pc.Queue) == 0 {
		return nil
	}

	qt := pc.Queue[0]
	pc.Queue = pc.Queue[1:]

	return &qt
}
