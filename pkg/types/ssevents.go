package types

import (
	"github.com/gofrs/uuid/v5"
)

type (
	SSEventType string
	SSENoData   struct{}
)

type SSEventBase struct {
	ID   uuid.UUID   `json:"id"`
	Type SSEventType `json:"type"`
	Data any         `json:"data"`
}

type SSEvent[T any] struct {
	ID   uuid.UUID   `json:"id"`
	Type SSEventType `json:"type"`
	Data T           `json:"data,omitempty"`
}

func (e SSEvent[T]) Base() SSEventBase {
	return SSEventBase{
		ID:   e.ID,
		Type: e.Type,
		Data: e.Data,
	}
}

type (
	SSEventPlayerPlayContext   = SSEvent[PlayContextDTO]
	SSEventPlayerPlaybackState = SSEvent[PlaybackStateDTO]
	SSEventPlayerStart         = SSEvent[PlayerState]
	SSEventPlayerPlayPause     = SSEvent[SSENoData]
)

const (
	SSEventTypePlayerPlayContext   SSEventType = "player.playcontext"
	SSEventTypePlayerPlaybackState SSEventType = "player.playbackstate"
	SSEventTypePlayerStart         SSEventType = "player.start"
	SSEventTypePlayerPlayPause     SSEventType = "player.playpause"
)

func newUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func SSEPlayerPlayContext(pc PlayContextDTO) SSEventPlayerPlayContext {
	return SSEventPlayerPlayContext{
		ID:   newUUID(),
		Type: SSEventTypePlayerPlayContext,
		Data: pc,
	}
}

func SSEPlayerPlaybackState(ps PlaybackStateDTO) SSEventPlayerPlaybackState {
	return SSEventPlayerPlaybackState{
		ID:   newUUID(),
		Type: SSEventTypePlayerPlaybackState,
		Data: ps,
	}
}

func SSEPlayerStart(pc PlayerState) SSEventPlayerStart {
	return SSEventPlayerStart{
		ID:   newUUID(),
		Type: SSEventTypePlayerStart,
		Data: pc,
	}
}

func SSEPlayerPlayPause() SSEventPlayerPlayPause {
	return SSEventPlayerPlayPause{
		ID:   newUUID(),
		Type: SSEventTypePlayerPlayPause,
	}
}
