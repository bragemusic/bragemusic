package musicbrainz

// ReleaseSearchResponse is the top-level structure returned from
// /ws/2/release/?query=...&fmt=json
type ReleaseSearchResponse struct {
	Created  string    `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []Release `json:"releases"`
}

type Release struct {
	ID                 string             `json:"id"`
	Score              int                `json:"score"`
	Title              string             `json:"title"`
	Status             string             `json:"status"`
	StatusID           string             `json:"status-id"`
	Packaging          string             `json:"packaging"`
	PackagingID        string             `json:"packaging-id"`
	ArtistCreditID     string             `json:"artist-credit-id"`
	Count              int                `json:"count"`
	TextRepresentation TextRepresentation `json:"text-representation"`
	ArtistCredit       []ArtistCredit     `json:"artist-credit"`
	ReleaseGroup       ReleaseGroup       `json:"release-group"`
	Date               MBDate             `json:"date"`
	Country            string             `json:"country"`
	ReleaseEvents      []ReleaseEvent     `json:"release-events"`
	Barcode            string             `json:"barcode"`
	LabelInfo          []LabelInfo        `json:"label-info"`
	TrackCount         int                `json:"track-count"`
	Media              []Media            `json:"media"`
	Relations          []Relation         `json:"relations"`
}

type TextRepresentation struct {
	Language string `json:"language"`
	Script   string `json:"script"`
}

type ArtistCredit struct {
	Name   string `json:"name"`
	Artist Artist `json:"artist"`
}

type Artist struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SortName string `json:"sort-name"`
}

type ReleaseGroup struct {
	ID            string `json:"id"`
	TypeID        string `json:"type-id"`
	PrimaryTypeID string `json:"primary-type-id"`
	Title         string `json:"title"`
	PrimaryType   string `json:"primary-type"`
}

type ReleaseEvent struct {
	Date MBDate `json:"date"`
	Area Area   `json:"area"`
}

type Area struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	SortName      string   `json:"sort-name"`
	ISO31661Codes []string `json:"iso-3166-1-codes"`
}

type LabelInfo struct {
	CatalogNumber string `json:"catalog-number"`
	Label         Label  `json:"label"`
}

type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// --- Media now includes Tracks ---
type Media struct {
	ID         string  `json:"id"`
	Format     string  `json:"format"`
	DiscCount  int     `json:"disc-count"`
	TrackCount int     `json:"track-count"`
	Tracks     []Track `json:"tracks"` // <-- new field
}

// --- Track details ---
type Track struct {
	ID        string    `json:"id"`
	Position  int       `json:"position"`
	Title     string    `json:"title"`
	Length    *int      `json:"length,omitempty"` // in milliseconds
	Recording Recording `json:"recording"`
}

// --- Recording info inside a Track ---
type Recording struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Length *int   `json:"length,omitempty"` // in milliseconds
}
