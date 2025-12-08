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
