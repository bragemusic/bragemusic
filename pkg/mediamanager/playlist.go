package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlaylist(ctx context.Context, p types.Playlist, userID uuid.UUID) error {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return m.berr.Unauthenticated(err)
	}

	if p.Name == "" {
		return m.berr.ParamMissing(nil, "name", types.EntityPlaylist.P(), types.ActionCreate.P())
	}

	p.Owner = user.ID

	if _, err := m.db.AddPlaylist(ctx, p, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, nil)
	}

	return nil
}

func (m MediaManager) CountPlaylists(ctx context.Context, userID uuid.UUID) (int, error) {
	cnt, err := m.db.CountPlaylists(ctx, userID)
	if err != nil {
		return 0, m.berr.DatabaseError(err, types.EntityPlaylist, nil)
	}

	return cnt, nil
}

func (m MediaManager) CountPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) (int, error) {
	plist, err := m.db.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return 0, m.berr.DatabaseError(err, types.EntityPlaylist, &playlistID)
	}

	if plist.Owner != userID {
		return 0, m.berr.ItemAccessDenied(nil, types.EntityPlaylist, playlistID)
	}

	cnt, err := m.db.CountPlaylistTracks(ctx, playlistID)
	if err != nil {
		return 0, m.berr.DatabaseError(err, types.EntityPlaylist, &playlistID)
	}

	return cnt, nil
}

func (m MediaManager) DeletePlaylist(ctx context.Context, playlistID, userID uuid.UUID) error {
	plist, err := m.db.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &playlistID)
	}

	if plist.Owner != userID {
		return m.berr.ItemAccessDenied(nil, types.EntityPlaylist, playlistID)
	}

	if err := m.db.DeletePlaylist(ctx, playlistID, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &playlistID)
	}

	return nil
}

func (m MediaManager) GetPlaylist(ctx context.Context, id uuid.UUID) (types.Playlist, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return types.Playlist{}, m.berr.Unauthenticated(err)
	}

	plist, err := m.db.GetPlaylist(ctx, id, user.ID)
	if err != nil {
		return types.Playlist{}, m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	return plist, nil
}

func (m MediaManager) ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, m.berr.Unauthenticated(err)
	}

	playlists, err := m.db.ListPlaylists(ctx, user.ID, includePublic, sortBy, sortOrder)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityPlaylist, nil)
	}

	return playlists, nil
}

func (m MediaManager) ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error) {
	plist, err := m.db.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityPlaylist, &playlistID)
	}

	if plist.Owner != userID {
		return nil, m.berr.ItemAccessDenied(nil, types.EntityPlaylist, playlistID)
	}

	tracks, err := m.db.ListPlaylistTracks(ctx, playlistID, userID)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityPlaylistTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist, userID uuid.UUID) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}
	defer tx.Rollback()

	existingPlist, err := tx.GetPlaylist(ctx, id, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	if existingPlist.Owner != userID {
		return m.berr.ItemAccessDenied(nil, types.EntityPlaylist, id)
	}

	data.ID = id
	data.Owner = existingPlist.Owner

	err = tx.UpdatePlaylist(ctx, data, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	err = tx.Commit()
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	return nil
}
