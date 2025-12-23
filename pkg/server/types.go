package server

import (
	"time"

	"github.com/bragemusic/core/pkg/types"
)

type HealthzStatus string

const (
	HealthzRunning HealthzStatus = "running"
)

type Healthz struct {
	Status HealthzStatus `json:"status"`
}

type Status struct {
	Application string        `json:"application"`
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
