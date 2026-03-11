package types

import (
	"errors"
	"math/rand/v2"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	PlayContextType string
	RepeatType      string
	TrackSource     string
)

const (
	PlayContextAlbum       PlayContextType = "album"
	PlayContextLikedTracks PlayContextType = "liked_tracks"
	PlayContextTopTracks   PlayContextType = "top_tracks"
	PlayContextPlaylist    PlayContextType = "playlist"

	RepeatOne RepeatType = "one"
	RepeatAll RepeatType = "all"
	RepeatOff RepeatType = "off"

	TrackSourceContext TrackSource = "context"
	TrackSourceQueue   TrackSource = "queue"
)

type PlayerState struct {
	Playback PlaybackState `json:"playback"`
	Context  PlayContext   `json:"context"`
}

type PlaybackState struct {
	Playing bool       `json:"playing"`
	Shuffle bool       `json:"shuffle"`
	Repeat  RepeatType `json:"repeat"`

	ProgressMS int64     `json:"progress"`
	UpdatedAt  time.Time `json:"updated_at"`

	TrackSource TrackSource `json:"track_source"`
	TrackIndex  int         `json:"track_index"`
}

func (p PlaybackState) TotalProgress() int64 {
	return p.ProgressMS + time.Now().Sub(p.UpdatedAt).Milliseconds()
}

type PlayContext struct {
	Type       PlayContextType `json:"type"`
	RefID      uuid.UUID       `json:"ref_id" ts_type:"string"`
	Tracks     []TrackDetailed `json:"tracks"`
	TrackOrder []int           `json:"track_order"`
	Queue      []TrackDetailed `json:"queue"`
}

func (pc *PlayContext) pullFromQueue() *TrackDetailed {
	if len(pc.Queue) == 0 {
		return nil
	}

	qt := pc.Queue[0]
	pc.Queue = pc.Queue[1:]

	return &qt
}

func (p *PlayerState) NextTrack() (contextUpdated, stop bool) {
	if p.Playback.TrackSource == TrackSourceQueue {
		queuedTrackExists := len(p.Context.Queue) > 0
		if queuedTrackExists {
			p.Context.pullFromQueue()
			contextUpdated = true
		}
	}

	queuedTrackExists := len(p.Context.Queue) > 0

	if !queuedTrackExists {
		var cidx int
		if p.Playback.Repeat == RepeatOne {
			cidx = p.Playback.TrackIndex
		} else {
			cidx = p.Playback.TrackIndex + 1
		}

		if cidx >= len(p.Context.Tracks) {
			if p.Playback.Repeat == RepeatAll {
				cidx = 0
			} else {
				return contextUpdated, true
			}
		}

		p.Playback.TrackSource = TrackSourceContext
		p.Playback.TrackIndex = cidx

		return contextUpdated, false
	} else {
		p.Playback.TrackSource = TrackSourceQueue
		return contextUpdated, false
	}
}

func (p *PlayerState) PreviousTrack() (stop bool) {
	var cidx int
	if p.Playback.Repeat == RepeatOne || p.Playback.TotalProgress() > 10000 {
		cidx = p.Playback.TrackIndex
	} else {
		cidx = p.Playback.TrackIndex - 1
	}

	if cidx < 0 {
		if p.Playback.Repeat == RepeatAll {
			cidx = len(p.Context.Tracks) - 1
		} else {
			return true
		}
	}

	p.Playback.TrackIndex = cidx

	return false
}

func (p *PlayerState) RebuildTrackOrder() {
	trackOrder := []int{}
	numberOfTracks := len(p.Context.Tracks)

	if numberOfTracks == 0 {
		p.Context.TrackOrder = []int{}
		return
	}

	trackIdx := p.Playback.TrackIndex
	if len(p.Context.TrackOrder) > p.Playback.TrackIndex {
		trackIdx = p.Context.TrackOrder[p.Playback.TrackIndex]
	}

	if !p.Playback.Shuffle {
		for i := range numberOfTracks {
			trackOrder = append(trackOrder, i)
		}
		p.Playback.TrackIndex = trackIdx
	} else {
		nums := make([]int, numberOfTracks)
		for i := range numberOfTracks {
			nums[i] = i
		}

		rand.Shuffle(numberOfTracks-1, func(i, j int) {
			nums[i+1], nums[j+1] = nums[j+1], nums[i+1]
		})

		for i := range numberOfTracks {
			if nums[i] == p.Playback.TrackIndex {
				nums[0], nums[i] = nums[i], nums[0]
				break
			}
		}

		trackOrder = nums
		p.Playback.TrackIndex = 0
	}

	p.Context.TrackOrder = trackOrder
}

func (p *PlayerState) CurrentTrack() (TrackDetailed, error) {
	switch p.Playback.TrackSource {
	case TrackSourceContext:
		if len(p.Context.TrackOrder) <= p.Playback.TrackIndex {
			return TrackDetailed{}, errors.New("track index oob")
		}
		idx := p.Context.TrackOrder[p.Playback.TrackIndex]
		if len(p.Context.Tracks) <= idx {
			return TrackDetailed{}, errors.New("track index oob")
		}
		return p.Context.Tracks[idx], nil

	case TrackSourceQueue:
		if len(p.Context.Queue) < 1 {
			return TrackDetailed{}, errors.New("track index oob")
		}
		return p.Context.Queue[0], nil

	default:
		return TrackDetailed{}, errors.New("unknown track source")
	}
}

type PlayContextDTO struct {
	PlaybackState   PlaybackState   `json:"playback_state"`
	Type            PlayContextType `json:"type"`
	RefID           uuid.UUID       `json:"ref_id" ts_type:"string"`
	Tracks          []uuid.UUID     `json:"tracks"`
	Queue           []uuid.UUID     `json:"queue"`
	CurrentTrackIdx int             `json:"current_track_idx"`
	CurrentTrack    *uuid.UUID      `json:"current_track"`
	TrackOrder      []int           `json:"track_order"`
}

type PlayContextOld struct {
	PlaybackState   PlaybackState   `json:"playback_state"`
	Type            PlayContextType `json:"type"`
	RefID           uuid.UUID       `json:"ref_id" ts_type:"string"`
	Tracks          []TrackDetailed `json:"tracks"`
	Queue           []TrackDetailed `json:"queue"`
	CurrentTrackIdx int             `json:"current_track_idx"`
	CurrentTrack    *TrackDetailed  `json:"current_track"`
	TrackOrder      []int           `json:"track_order"`
}

func (pc *PlayContextOld) PullFromQueue() *TrackDetailed {
	if len(pc.Queue) == 0 {
		return nil
	}

	qt := pc.Queue[0]
	pc.Queue = pc.Queue[1:]

	return &qt
}

// func BuildPlayContext(dto PlayContextDTO, tracks []TrackDetailed, queue []TrackDetailed) PlayContext {
// 	pc := PlayContext{
// 		Type:            dto.Type,
// 		RefID:           dto.RefID,
// 		Tracks:          tracks,
// 		Queue:           queue,
// 		CurrentTrackIdx: dto.CurrentTrackIdx,
// 		TrackOrder:      dto.TrackOrder,
// 		PlaybackState:   dto.PlaybackState,
// 	}

// 	if dto.CurrentTrack != nil {
// 		for _, t := range append(tracks, queue...) {
// 			if t.ID == *dto.CurrentTrack {
// 				pc.CurrentTrack = &t
// 				break
// 			}
// 		}
// 	}

// 	return pc
// }
