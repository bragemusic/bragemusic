package types

import (
	"encoding/json"

	"github.com/gofrs/uuid/v5"
)

type (
	SSEventType string
	SSENoData   struct{}
)

// type SSEventBase struct {
// 	ID   uuid.UUID   `json:"id"`
// 	Type SSEventType `json:"type"`
// 	Data any         `json:"data"`
// }

type SSEvent struct {
	ID   uuid.UUID       `json:"id"`
	Type SSEventType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

// func (e SSEvent[T]) Base() SSEventBase {
// 	return SSEventBase{
// 		ID:   e.ID,
// 		Type: e.Type,
// 		Data: e.Data,
// 	}
// }

// type (
// 	SSEventClientConnected     = SSEvent[Device]
// 	SSEventClientDisconnected  = SSEvent[Device]
// 	SSEventPlayerPlayContext   = SSEvent[PlayContextDTO]
// 	SSEventPlayerPlaybackState = SSEvent[PlaybackStateDTO]
// 	SSEventPlayerStart         = SSEvent[PlayerState]
// 	SSEventPlayerPlayPause     = SSEvent[SSENoData]
// )

const (
	SSEventTypeDeviceConnected     SSEventType = "device.connected"
	SSEventTypeDeviceUpdated       SSEventType = "device.updated"
	SSEventTypeDeviceDisconnected  SSEventType = "device.disconnected"
	SSEventTypeDevicePlayContext   SSEventType = "device.playcontext"
	SSEventTypeDevicePlaybackState SSEventType = "device.playbackstate"
	SSEventTypePlayerPlayContext   SSEventType = "player.playcontext"
	SSEventTypePlayerPlaybackState SSEventType = "player.playbackstate"
	SSEventTypePlayerSetRepeat     SSEventType = "player.setrepeat"
	SSEventTypePlayerSetShuffle    SSEventType = "player.setshuffle"
	SSEventTypePlayerSetState      SSEventType = "player.setstate"
	SSEventTypePlayerStop          SSEventType = "player.stop"
	SSEventTypePlayerNextTrack     SSEventType = "player.nexttrack"
	SSEventTypePlayerPlayPause     SSEventType = "player.playpause"
	SSEventTypePlayerPreviousTrack SSEventType = "player.previoustrack"
)

func DecodeEventData[T any](evt SSEvent) (T, error) {
	var v T
	err := json.Unmarshal(evt.Data, &v)
	return v, err
}

func newUUID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func newSSEvent(eventType SSEventType, data any) SSEvent {
	b, err := json.Marshal(data)
	if err != nil {
		panic("generating event of type " + string(eventType) + ". ERROR: " + err.Error())
	}

	return SSEvent{
		ID:   newUUID(),
		Type: eventType,
		Data: b,
	}
}

func SSEDeviceConnected(d Device) SSEvent {
	return newSSEvent(SSEventTypeDeviceConnected, d)
}

func SSEDeviceUpdated(d Device) SSEvent {
	return newSSEvent(SSEventTypeDeviceUpdated, d)
}

func SSEDeviceDisconnected(d Device) SSEvent {
	return newSSEvent(SSEventTypeDeviceDisconnected, d)
}

func SSEDevicePlayContext(pc PlayContextDTO) SSEvent {
	return newSSEvent(SSEventTypeDevicePlayContext, pc)
}

func SSEDevicePlaybackState(ps PlaybackStateDTO) SSEvent {
	return newSSEvent(SSEventTypeDevicePlaybackState, ps)
}

func SSEPlayerPlayContext(pc PlayContextDTO) SSEvent {
	return newSSEvent(SSEventTypePlayerPlayContext, pc)
}

func SSEPlayerPlaybackState(ps PlaybackStateDTO) SSEvent {
	return newSSEvent(SSEventTypePlayerPlaybackState, ps)
}

func SSEPlayerSetRepeat(rt RepeatType) SSEvent {
	return newSSEvent(SSEventTypePlayerSetRepeat, rt)
}

func SSEPlayerSetShuffle(a bool) SSEvent {
	return newSSEvent(SSEventTypePlayerSetShuffle, a)
}

func SSEPlayerSetState(pc PlayerState) SSEvent {
	return newSSEvent(SSEventTypePlayerSetState, pc)
}

func SSEPlayerStop() SSEvent {
	return newSSEvent(SSEventTypePlayerStop, nil)
}

func SSEPlayerNextTrack() SSEvent {
	return newSSEvent(SSEventTypePlayerNextTrack, nil)
}

func SSEPlayerPlayPause() SSEvent {
	return newSSEvent(SSEventTypePlayerPlayPause, nil)
}

func SSEPlayerPreviousTrack() SSEvent {
	return newSSEvent(SSEventTypePlayerPreviousTrack, nil)
}
