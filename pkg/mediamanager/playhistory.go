package mediamanager

import (
	"context"
	"fmt"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlayCount(ctx context.Context, trackID, userID string) error {
	trackUUID, err := uuid.FromString(trackID)
	if err != nil {
		return fmt.Errorf("could not parse trackID: %s", err.Error())
	}

	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return fmt.Errorf("could not parse userID: %s", err.Error())
	}

	if _, err = m.db.AddPlayHistory(ctx, trackUUID, userUUID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlayHistoryItem, nil)
	}

	return nil
}
