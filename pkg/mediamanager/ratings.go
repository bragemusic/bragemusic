package mediamanager

import (
	"context"
	"database/sql"
	"errors"

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

func (m MediaManager) GetRating(ctx context.Context, id uuid.UUID) (types.Rating, error) {
	rating, err := m.db.GetRating(ctx, id)
	if err != nil {
		return types.Rating{}, m.berr.DatabaseError(err, types.EntityRating, &id)
	}

	return rating, nil
}

func (m MediaManager) GetTrackRatings(ctx context.Context, trackID uuid.UUID) ([]types.Rating, error) {
	ratings, err := m.db.GetTrackRatings(ctx, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, m.berr.DatabaseError(err, types.EntityRating, nil)
	}

	return ratings, nil
}
