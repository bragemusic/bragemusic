package mediamanager

import (
	"context"
	"errors"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlaylist(ctx context.Context, p types.Playlist) error {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return err
	}

	if p.Name == "" {
		return errors.New("cannot create playlist, no name was provided")
	}

	p.Owner = user.ID

	if _, err := m.db.AddPlaylist(ctx, p); err != nil {
		return err
	}

	return nil
}

func (m MediaManager) GetPlaylist(ctx context.Context, id uuid.UUID) (types.Playlist, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return types.Playlist{}, err
	}

	plist, err := m.db.GetPlaylist(ctx, id, user.ID)
	if err != nil {
		return types.Playlist{}, err
	}

	return plist, nil
}

func (m MediaManager) ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	playlists, err := m.db.ListPlaylists(ctx, user.ID, includePublic, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	return playlists, nil
}
