package types

import "time"

type PlayHistory struct {
	ID       string    `db:"id" json:"id"`
	UserID   string    `db:"user_id" json:"user_id"`
	TrackID  string    `db:"track_id" json:"track_id"`
	PlayedAt time.Time `db:"played_at" json:"played_at"`
}
