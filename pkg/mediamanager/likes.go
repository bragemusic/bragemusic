package mediamanager

import (
	"context"
	"errors"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddLike(ctx context.Context, trackID, userID uuid.UUID) error {
	hasLike, err := m.db.HasLike(ctx, trackID, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityLike, nil)
	}

	if hasLike {
		return m.berr.EntityAlreadyExists(errors.New("user cannot like an already liked track"), types.EntityLike)
	}

	l := types.Like{
		TrackID: trackID,
		Owner:   userID,
	}

	if _, err := m.db.AddLike(ctx, l); err != nil {
		return m.berr.DatabaseError(err, types.EntityLike, nil)
	}

	return nil
}

func (m MediaManager) GetLike(ctx context.Context, id, userID uuid.UUID) (types.Like, error) {
	like, err := m.db.GetLike(ctx, id, userID)
	if err != nil {
		return types.Like{}, m.berr.DatabaseError(err, types.EntityLike, &id)
	}

	return like, nil
}

func (m MediaManager) ListLikedTracksDetailed(ctx context.Context, userID uuid.UUID) (tracks []types.TrackDetailed, err error) {
	tracks, err = m.db.ListLikedTracksDetailed(ctx, userID)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) RemoveLike(ctx context.Context, trackID, userID uuid.UUID) error {
	hasLike, err := m.db.HasLike(ctx, trackID, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityLike, nil)
	}

	if !hasLike {
		return m.berr.EntityDoNotExist(errors.New("user cannot remove an unliked track"), types.EntityLike)
	}

	id, err := m.db.GetLikeID(ctx, trackID, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityLike, nil)
	}

	if err := m.db.DeleteLike(ctx, id, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityLike, &id)
	}

	return nil
}
