package audioplayer

import (
	"github.com/bragemusic/core/pkg/types"
)

type PlayContextType string

const (
	PlayContextAlbum     PlayContextType = "album"
	PlayContextTopTracks PlayContextType = "top_tracks"
)

type PlayContext struct {
	Type            PlayContextType       `json:"type"`
	RefID           string                `json:"ref_id"`
	Tracks          []types.TrackEnhanced `json:"tracks"`
	Queue           []types.TrackEnhanced `json:"queue"`
	CurrentTrackIdx int                   `json:"current_track_idx"`
}
