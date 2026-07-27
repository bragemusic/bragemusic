package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlaylist(ctx context.Context, p types.PlaylistBase, userID uuid.UUID) error {
	if p.Name == "" {
		return m.berr.ParamMissing(nil, "name", types.EntityPlaylist.P(), types.ActionCreate.P())
	}

	pl := types.Playlist{
		PlaylistBase: p,
		Owner:        userID,
		Type:         types.PlaylistTypeStandard,
	}

	if _, err := m.db.AddPlaylist(ctx, pl, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, nil)
	}

	return nil
}

func (m MediaManager) AddSmartPlaylist(ctx context.Context, p types.PlaylistBase, filter types.TrackFilter, userID uuid.UUID) error {
	if p.Name == "" {
		return m.berr.ParamMissing(nil, "name", types.EntityPlaylist.P(), types.ActionCreate.P())
	}

	c := types.SmartPlaylistContent{
		MoodAggressive: filter.Mood.Aggressive,
		MoodCalm:       filter.Mood.Calm,
		MoodHappy:      filter.Mood.Happy,
		MoodSad:        filter.Mood.Sad,
		Artists:        filter.Artists,
	}

	if filter.BPM != nil {
		c.BpmLower = &filter.BPM.Lower
		c.BpmUpper = &filter.BPM.Upper
	}

	pl := types.SmartPlaylist{
		Playlist: types.Playlist{
			PlaylistBase: p,
			Owner:        userID,
		},
		Content: c,
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.AddSmartPlaylist(ctx, pl, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, nil)
	}

	err = tx.Commit()
	if err != nil {
		return err
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

func (m MediaManager) GetPlaylist(ctx context.Context, id, userID uuid.UUID) (types.Playlist, error) {
	plist, err := m.db.GetPlaylist(ctx, id, userID)
	if err != nil {
		return types.Playlist{}, m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	return plist, nil
}

func (m MediaManager) GetSmartPlaylistContent(ctx context.Context, id, userID uuid.UUID) (types.SmartPlaylistContent, error) {
	content, err := m.db.GetSmartPlaylistContent(ctx, id, userID)
	if err != nil {
		return types.SmartPlaylistContent{}, m.berr.DatabaseError(err, types.EntitySmartPlaylistContent, &id)
	}

	return content, nil
}

func (m MediaManager) ListPlaylists(ctx context.Context, userID uuid.UUID, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error) {
	playlists, err := m.db.ListPlaylists(ctx, userID, includePublic, sortBy, sortOrder)
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

func (m MediaManager) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.PlaylistBase, userID uuid.UUID) error {
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

	existingPlist.PlaylistBase = data

	err = tx.UpdatePlaylist(ctx, existingPlist, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	err = tx.Commit()
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	return nil
}
