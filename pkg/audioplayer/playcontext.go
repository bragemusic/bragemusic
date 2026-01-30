package audioplayer

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type PlayContextType string

const (
	PlayContextAlbum     PlayContextType = "album"
	PlayContextTopTracks PlayContextType = "top_tracks"
	PlayContextPlaylist  PlayContextType = "playlist"
)

type PlayContext struct {
	Type            PlayContextType       `json:"type"`
	RefID           uuid.UUID             `json:"ref_id" ts_type:"string"`
	Tracks          []types.TrackDetailed `json:"tracks"`
	Queue           []types.TrackDetailed `json:"queue"`
	CurrentTrackIdx int                   `json:"current_track_idx"`
	CurrentTrack    *types.TrackDetailed  `json:"current_track"`
}
