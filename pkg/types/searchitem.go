package types

import "github.com/gofrs/uuid/v5"

type SearchItem struct {
	Name     string     `db:"name" json:"name"`
	HtmlName string     `db:"html_name" json:"html_name"`
	ID       uuid.UUID  `db:"id" json:"id"`
	Type     EntityType `db:"type" json:"type"`
	LinkID   uuid.UUID  `db:"link_id" json:"link_id"`
	LinkType EntityType `db:"link_type" json:"link_type"`
}
