package musicbrainz

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SearchResults struct {
	ID     string `json:"id"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Score  int    `json:"score"`
}

type MBDate struct {
	time.Time
}

// UnmarshalJSON parses a MusicBrainz-style date.
func (d *MBDate) UnmarshalJSON(b []byte) error {
	// Remove quotes
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}

	// Try progressively more detailed layouts
	layouts := []string{
		"2006-01-02",
		"2006-01",
		"2006",
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			d.Time = t
			return nil
		}
	}

	return fmt.Errorf("invalid MusicBrainz date: %q", s)
}

// MarshalJSON ensures consistent ISO-style output if you re-encode it.
func (d MBDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(d.Format("2006-01-02"))
}
