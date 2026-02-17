package mediamanager

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) ListEntityEvents(ctx context.Context) ([]types.EntityEvent, error) {
	events, err := m.db.ListEntityEvents(ctx, time.Unix(0, 0))
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityEntityEvent, nil)
	}

	return events, nil
}
