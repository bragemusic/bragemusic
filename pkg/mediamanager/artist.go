package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) CreateArtist(ctx context.Context, artistData types.ArtistBase, userID uuid.UUID) error {
	artist := types.Artist{
		ArtistBase: artistData,
	}

	_, err := m.db.AddArtist(ctx, artist, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityArtist, nil)
	}

	return nil
}

func (m MediaManager) GetArtist(ctx context.Context, artistID uuid.UUID) (types.Artist, error) {
	artist, err := m.db.GetArtistFromID(ctx, artistID)
	if err != nil {
		return types.Artist{}, m.berr.DatabaseError(err, types.EntityArtist, &artistID)
	}

	return artist, nil
}

func (m MediaManager) ListArtists(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.ArtistDetailed, error) {
	artists, err := m.db.ListArtists(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityArtist, nil)
	}

	return artists, nil
}

func (m MediaManager) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist, userID uuid.UUID) error {
	artistData.ID = artistID
	err := m.db.UpdateArtist(ctx, artistData, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityArtist, &artistID)
	}

	return nil
}

func (m MediaManager) CountArtists(ctx context.Context) (int, error) {
	cnt, err := m.db.CountArtists(ctx)
	if err != nil {
		return 0, m.berr.DatabaseError(err, types.EntityArtist, nil)
	}

	return cnt, nil
}
