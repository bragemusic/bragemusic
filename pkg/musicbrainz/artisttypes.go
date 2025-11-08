package musicbrainz

// ArtistResponse represents the full artist object returned from
// /ws/2/artist/{id}?inc=aliases+tags+genres&fmt=json
type ArtistResponse struct {
	ID             string      `json:"id"`
	Type           string      `json:"type"`
	TypeID         string      `json:"type-id"`
	Name           string      `json:"name"`
	SortName       string      `json:"sort-name"`
	Country        string      `json:"country"`
	Gender         *string     `json:"gender"`
	GenderID       *string     `json:"gender-id"`
	Disambiguation string      `json:"disambiguation"`
	LifeSpan       LifeSpan    `json:"life-span"`
	Area           AreaSimple  `json:"area"`
	BeginArea      *AreaSimple `json:"begin-area"`
	EndArea        *AreaSimple `json:"end-area"`
	ISNIs          []string    `json:"isnis"`
	IPIs           []string    `json:"ipis"`
	Aliases        []Alias     `json:"aliases"`
	Tags           []Tag       `json:"tags"`
	Genres         []Genre     `json:"genres"`
	Relations      []Relation  `json:"relations"`
}

type LifeSpan struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
	Ended bool   `json:"ended"`
}

type AreaSimple struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	SortName       string   `json:"sort-name"`
	Disambiguation string   `json:"disambiguation"`
	Type           *string  `json:"type"`
	TypeID         *string  `json:"type-id"`
	ISO31661Codes  []string `json:"iso-3166-1-codes"`
}

type Alias struct {
	Name     string  `json:"name"`
	SortName string  `json:"sort-name"`
	Locale   *string `json:"locale"`
	Type     string  `json:"type"`
	TypeID   string  `json:"type-id"`
	Primary  *bool   `json:"primary"`
	Begin    *string `json:"begin"`
	End      *string `json:"end"`
	Ended    bool    `json:"ended"`
}

type Tag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Genre struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Disambiguation string `json:"disambiguation"`
	Count          int    `json:"count"`
}

type Relation struct {
	Type         string            `json:"type"`
	TypeID       string            `json:"type-id"`
	Direction    string            `json:"direction"`
	TargetType   string            `json:"target-type"`
	OrderingKey  int               `json:"ordering-key,omitempty"`
	AttributeIDs map[string]string `json:"attribute-ids,omitempty"`
	Attributes   []string          `json:"attributes,omitempty"`
	URL          struct {
		ID       string `json:"id"`
		Resource string `json:"resource"`
	} `json:"url"`
}
