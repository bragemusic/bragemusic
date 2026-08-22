package mediamanager

import (
	"context"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) GetMediaFile(ctx context.Context, mediafileID uuid.UUID) (types.MediaFile, error) {
	mf, err := m.db.GetMediaFile(ctx, mediafileID)
	if err != nil {
		return types.MediaFile{}, m.berr.DatabaseError(err, types.EntityMediaFile, &mediafileID)
	}

	return mf, nil
}
