package database

import (
	"context"
	"strings"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/jmoiron/sqlx"
)

func (d Database) AddSearchItem(ctx context.Context, si types.SearchItem) error {
	query := `
        INSERT INTO search_items(
            name, id, type, link_id, link_type
        )
        VALUES (?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		si.Name,
		si.ID,
		si.Type,
		si.LinkID,
		si.LinkType,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) DeleteAllSearchItems(ctx context.Context) error {
	query := `
        DELETE FROM search_items;
   `

	_, err := d.ext.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) SearchFull(ctx context.Context, searchTerm string, limit int) (results []types.SearchItem, err error) {
	if searchTerm == "" {
		return nil, nil
	}

	term := strings.TrimSpace(searchTerm)
	term = strings.ReplaceAll(term, `"`, "")
	term = strings.ReplaceAll(term, `'`, "")

	query := `
        SELECT
          name,
          highlight(search_items, 0, '<b><u>', '</b></u>') html_name,
          id,
          type,
          link_id,
          link_type
        FROM search_items
        WHERE search_items MATCH ?
        ORDER BY rank
        LIMIT ?;
    `
	err = sqlx.SelectContext(ctx, d.ext, &results, query, term+"*", limit)
	if err != nil {
		return nil, err
	}

	return
}
