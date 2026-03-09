package types

import (
	"github.com/gofrs/uuid/v5"
)

type (
	PlayContextType string
	RepeatType      string
)

const (
	PlayContextAlbum       PlayContextType = "album"
	PlayContextLikedTracks PlayContextType = "liked_tracks"
	PlayContextTopTracks   PlayContextType = "top_tracks"
	PlayContextPlaylist    PlayContextType = "playlist"

	RepeatOne RepeatType = "one"
	RepeatAll RepeatType = "all"
	RepeatOff RepeatType = "off"
)

type PlayContext struct {
	Type            PlayContextType `json:"type"`
	RefID           uuid.UUID       `json:"ref_id" ts_type:"string"`
	Tracks          []TrackDetailed `json:"tracks"`
	Queue           []TrackDetailed `json:"queue"`
	CurrentTrackIdx int             `json:"current_track_idx"`
	CurrentTrack    *TrackDetailed  `json:"current_track"`
	Shuffle         bool            `json:"shuffle"`
	Repeat          RepeatType      `json:"repeat"`
	TrackOrder      []int           `json:"track_order"`
}

func (pc *PlayContext) PullFromQueue() *TrackDetailed {
	if len(pc.Queue) == 0 {
		return nil
	}

	qt := pc.Queue[0]
	pc.Queue = pc.Queue[1:]

	return &qt
}
