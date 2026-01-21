package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) SearchFull(ctx context.Context, searchTerm string) (si []types.SearchItem, err error) {
	if searchTerm == "" {
		return nil, nil
	}

	res, err := m.db.SearchFull(ctx, searchTerm, 10)
	if err != nil {
		return nil, err
	}

	return res, nil
}
