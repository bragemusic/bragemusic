package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) GetArtist(ctx context.Context, artistID string) (types.Artist, error) {
	artist, err := m.db.GetArtistFromID(ctx, artistID)
	if err != nil {
		return types.Artist{}, err
	}

	return artist, nil
}

func (m MediaManager) ListArtists(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Artist, error) {
	artists, err := m.db.ListArtists(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (m MediaManager) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	artistData.ID = artistID
	err := m.db.UpdateArtist(ctx, artistData)
	if err != nil {
		return err
	}

	return nil
}
