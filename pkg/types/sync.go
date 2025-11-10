package types

import "time"

type SyncState struct {
	Time    time.Time `json:"time"`
	Artists []string  `json:"artists"`
	Albums  []string  `json:"albums"`
	Tracks  []string  `json:"tracks"`
}
