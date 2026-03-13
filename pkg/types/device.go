package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	DeviceType      string
	DeviceInterface string
)

const (
	DeviceTypeStreaming DeviceType = "streaming"
	DeviceTypeSync      DeviceType = "sync"

	DeviceInterfaceDaemon   DeviceInterface = "daemon"
	DeviceInterfaceDesktop  DeviceInterface = "desktop"
	DeviceInterfaceTerminal DeviceInterface = "terminal"
	DeviceInterfaceWeb      DeviceInterface = "web"
)

// TODO:UNIQUE(token_id) (maybe, might not be necessary. Should not be able to disable a client, but a token. A token should be able to be used in multiple places)
type DeviceBase struct {
	Name             string          `db:"name" json:"name" `
	Type             DeviceType      `db:"type" json:"type"`
	Interface        DeviceInterface `db:"interface" json:"interface"`
	SupportsPlayback bool            `db:"supports_playback" json:"supports_playback"`
	Platform         string          `db:"platform" json:"platform"`
	Version          string          `db:"version" json:"version"`
}

type ReqDevicesRegister struct {
	ID *uuid.UUID `json:"id"`
	DeviceBase
}

type Device struct {
	DeviceBase
	ID        uuid.UUID `db:"id" json:"id" ts_type:"string"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	LastIP    string    `db:"last_ip" json:"last_ip"`
	LastSeen  time.Time `db:"last_seen" json:"last_seen"`
	Active    bool      `db:"-" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type DeviceDetailed struct {
	Device
	PlayerState *PlayerStateDTO `json:"player_state"`
}

type DeviceToken struct {
	DeviceID  uuid.UUID `db:"device_id"`
	TokenID   uuid.UUID `db:"token_id"`
	CreatedAt time.Time `db:"created_at"`
}
