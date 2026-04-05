package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type HealthzStatus string

const (
	HealthzRunning     HealthzStatus = "running"
	HealthzUnavailable HealthzStatus = "unavailable"
	HealthzNoAuth      HealthzStatus = "no_auth"
)

type Healthz struct {
	Status HealthzStatus `json:"status"`
}

type ServerInfo struct {
	Application string        `json:"application"`
	ID          uuid.UUID     `json:"id" ts_type:"string"`
	Status      HealthzStatus `json:"status"`
}

type ServerApiInfo struct {
	ServerInfo
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SyncReq struct {
	ChangesSince time.Time `json:"changes_since"`
}

type SyncPlayHistoryReq struct {
	ChangesSince       time.Time     `json:"changes_since"`
	UpdatedClientItems []PlayHistory `json:"updated_client_items"`
}

type PlaylistTrackReq struct {
	AlbumID uuid.UUID `json:"album_id"`
	TrackID uuid.UUID `json:"track_id"`
}

type RatingReq struct {
	Value int `json:"value"`
}

type ListPayload[T any] struct {
	Items []T `json:"items,omitempty"`
	Count int `json:"count"`
}

type LoginReq struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	LongLivedToken bool   `json:"long_lived_token"`
}

type LoginResp struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

type RespDevicesRegister struct {
	DeviceID uuid.UUID `json:"id"`
}

type SetDevicePlayerRepeatReq struct {
	Type RepeatType `json:"type" description:"Type of repeat"`
}

type SetDevicePlayerShuffleReq struct {
	Active bool `json:"active" description:"Shuffle is active"`
}

type CreateUserReq struct {
	Email    string     `json:"email"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Roles    []UserRole `json:"roles"`
}

type UpdateUserReq struct {
	Email    string     `json:"email"`
	Username string     `json:"username"`
	Password *string    `json:"password"`
	Roles    []UserRole `json:"roles"`
}
