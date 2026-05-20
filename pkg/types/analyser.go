package types

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type TrackAnalysisResults struct {
	BPM           int     `db:"bpm" json:"bpm"`
	Key           string  `db:"key" json:"key"`
	KeyScale      string  `db:"key_scale" json:"key_scale"`
	KeyConfidence float64 `db:"key_confidence" json:"key_confidence"`
	Loudness      int     `db:"loudness" json:"loudness"`
	Energy        float64 `db:"energy" json:"energy"`
	Danceability  float64 `db:"danceability" json:"danceability"`
	MoodHappy     float64 `db:"mood_happy" json:"mood_happy"`
	MoodSad       float64 `db:"mood_sad" json:"mood_sad"`
	MoodAggresive float64 `db:"mood_aggresive" json:"mood_aggresive"`
	MoodCalm      float64 `db:"mood_calm" json:"mood_calm"`
}

type TrackAnalysis struct {
	ID uuid.UUID `db:"id" json:"id" ts_type:"string"`
	TrackAnalysisResults
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
