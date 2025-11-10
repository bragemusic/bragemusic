package server

import "time"

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
