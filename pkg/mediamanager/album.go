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

func (m MediaManager) GetAlbumArtist(ctx context.Context, albumID, artistID, role string) (types.AlbumArtist, error) {
	albumUID, err := uuid.FromString(albumID)
	if err != nil {
		return types.AlbumArtist{}, err
	}

	artistUID, err := uuid.FromString(artistID)
	if err != nil {
		return types.AlbumArtist{}, err
	}

	albumArtist, err := m.db.GetAlbumArtist(ctx, albumUID, artistUID, types.ArtistRole(role))
	if err != nil {
		return types.AlbumArtist{}, err
	}

	return albumArtist, nil
}

func (m MediaManager) GetAlbumTrack(ctx context.Context, albumID uuid.UUID, discNumber, trackNumber int) (types.AlbumTrack, error) {
	albumArtist, err := m.db.GetAlbumTrack(ctx, albumID, discNumber, trackNumber)
	if err != nil {
		return types.AlbumTrack{}, err
	}

	return albumArtist, nil
}
