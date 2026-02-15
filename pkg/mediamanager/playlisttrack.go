package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID, userID uuid.UUID) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylistTrack, nil)
	}
	defer tx.Rollback()

	plist, err := tx.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &playlistID)
	}

	if plist.Owner != userID {
		return m.berr.ItemAccessDenied(nil, types.EntityPlaylist, playlistID)
	}

	albumTrack, err := tx.GetAlbumTrackFromAlbumAndTrack(ctx, albumID, trackID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityAlbumTrack, nil)
	}

	plistTrack := types.PlaylistTrack{
		PlaylistID:   playlistID,
		AlbumTrackID: albumTrack.ID,
	}

	if _, err := tx.AddPlaylistTrack(ctx, plistTrack); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylistTrack, nil)
	}

	return tx.Commit()
}

func (m MediaManager) DeletePlaylistTrack(ctx context.Context, id, userID uuid.UUID) error {
	pt, err := m.GetPlaylistTrack(ctx, id, userID)
	if err != nil {
		return err
	}

	plist, err := m.GetPlaylist(ctx, pt.PlaylistID)
	if err != nil {
		return err
	}

	if plist.Owner != userID {
		return m.berr.ItemAccessDenied(nil, types.EntityPlaylist, pt.PlaylistID)
	}

	if err = m.db.DeletePlaylistTrack(ctx, pt.ID, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylistTrack, &id)
	}

	return nil
}

func (m MediaManager) GetPlaylistTrack(ctx context.Context, id, userID uuid.UUID) (types.PlaylistTrack, error) {
	pt, err := m.db.GetPlaylistTrack(ctx, id)
	if err != nil {
		return types.PlaylistTrack{}, m.berr.DatabaseError(err, types.EntityPlaylistTrack, &id)
	}

	plist, err := m.GetPlaylist(ctx, pt.PlaylistID)
	if err != nil {
		return types.PlaylistTrack{}, err
	}

	if plist.Owner != userID {
		return types.PlaylistTrack{}, m.berr.ItemAccessDenied(nil, types.EntityPlaylist, plist.ID)
	}

	return pt, nil
}
