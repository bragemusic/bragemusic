package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlayCount(ctx context.Context, trackID, userID uuid.UUID) error {
	if _, err := m.db.AddPlayHistory(ctx, trackID, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlayHistoryItem, nil)
	}

	m.log.DebugContext(ctx, "added play count", "track_id", trackID, "user", userID)

	return nil
}
