package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) GetAlbum(ctx context.Context, albumID string) (types.Album, error) {
	album, err := m.db.GetAlbumFromID(ctx, albumID)
	if err != nil {
		return types.Album{}, err
	}

	return album, nil
}

func (m MediaManager) GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (types.AlbumDetailed, error) {
	album, err := m.db.GetAlbumDetailed(ctx, albumID)
	if err != nil {
		return types.AlbumDetailed{}, err
	}

	return album, nil
}

func (m MediaManager) ListAlbumsByArtist(ctx context.Context, artistID string, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error) {
	albums, err := m.db.ListAlbumsByArtist(ctx, artistID, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	return albums, nil
}
