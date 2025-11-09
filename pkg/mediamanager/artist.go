package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) GetArtist(ctx context.Context, artistID string) (types.Artist, error) {
	artist, err := m.db.GetArtistFromID(ctx, artistID)
	if err != nil {
		return types.Artist{}, err
	}

	return artist, nil
}

func (m MediaManager) ListArtists(ctx context.Context) ([]types.Artist, error) {
	artists, err := m.db.ListArtists(ctx)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
