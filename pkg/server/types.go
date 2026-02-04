package server

import (
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type HealthzStatus string

const (
	HealthzRunning     HealthzStatus = "running"
	HealthzUnavailable HealthzStatus = "unavailable"
)

type Healthz struct {
	Status HealthzStatus `json:"status"`
}

type Response struct {
	Payload any `json:"payload,omitempty"`
	Status  int `json:"-"`
}

type Status struct {
	Application string        `json:"application"`
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Status      HealthzStatus `json:"status"`
}

type SyncReq struct {
	ChangesSince time.Time `json:"changes_since"`
}

type SyncPlayHistoryReq struct {
	ChangesSince       time.Time           `json:"changes_since"`
	UpdatedClientItems []types.PlayHistory `json:"updated_client_items"`
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

type PlaylistTrackReq struct {
	AlbumID uuid.UUID `json:"album_id"`
	TrackID uuid.UUID `json:"track_id"`
}
