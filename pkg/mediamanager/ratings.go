package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) RateTrack(ctx context.Context, trackID, userID uuid.UUID, value int) error {
	rateID, found, err := m.db.GetRatingID(ctx, trackID, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityRating, nil)
	}

	if found {
		if err = m.db.UpdateRating(ctx, rateID, value); err != nil {
			return m.berr.DatabaseError(err, types.EntityRating, &rateID)
		}
		return nil
	}

	r := types.Rating{
		TrackID: trackID,
		Rating:  value,
		Owner:   userID,
	}

	if _, err = m.db.AddRating(ctx, r); err != nil {
		return m.berr.DatabaseError(err, types.EntityRating, nil)
	}

	return nil
}
