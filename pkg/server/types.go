package server

type HealthzStatus string

const (
	HealthzRunning HealthzStatus = "running"
)

type Healthz struct {
	Application string        `json:"application"`
	Version     string        `json:"version"`
	Status      HealthzStatus `json:"status"`
}
