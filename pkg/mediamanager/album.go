package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) GetAlbum(ctx context.Context, albumID string) (types.Album, error) {
	album, err := m.db.GetAlbumFromID(ctx, albumID)
	if err != nil {
		return types.Album{}, err
	}

	return album, nil
}

func (m MediaManager) ListAlbumsByArtist(ctx context.Context, artistID string) ([]types.Album, error) {
	albums, err := m.db.ListAlbumsByArtist(ctx, artistID)
	if err != nil {
		return nil, err
	}

	return albums, nil
}
