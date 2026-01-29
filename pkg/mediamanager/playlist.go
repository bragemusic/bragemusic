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

func (m MediaManager) CountPlaylists(ctx context.Context, userID uuid.UUID) (int, error) {
	return m.db.CountPlaylists(ctx, userID)
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

func (m MediaManager) ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error) {
	plist, err := m.db.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return nil, err
	}

	if plist.Owner != userID {
		return nil, errors.New("user is not the owner of the selected playlist")
	}

	tracks, err := m.db.ListPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (m MediaManager) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist, userID uuid.UUID) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existingPlist, err := tx.GetPlaylist(ctx, id, userID)
	if err != nil {
		return err
	}

	data.ID = id
	data.Owner = existingPlist.Owner

	err = tx.UpdatePlaylist(ctx, data)
	if err != nil {
		return err
	}

	return tx.Commit()
}
