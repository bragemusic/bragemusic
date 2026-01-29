package mediamanager

import (
	"context"
	"errors"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID, userID uuid.UUID) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	plist, err := tx.GetPlaylist(ctx, playlistID, userID)
	if err != nil {
		return err
	}

	if plist.Owner != userID {
		return errors.New("user is not the owner of the selected playlist")
	}

	albumTrack, err := tx.GetAlbumTrackFromAlbumAndTrack(ctx, albumID, trackID)
	if err != nil {
		return err
	}

	plistTrack := types.PlaylistTrack{
		PlaylistID:   playlistID,
		AlbumTrackID: albumTrack.ID,
	}

	if _, err := tx.AddPlaylistTrack(ctx, plistTrack); err != nil {
		return err
	}

	return tx.Commit()
}

func (m MediaManager) GetPlaylistTrack(ctx context.Context, id, userID uuid.UUID) (types.PlaylistTrack, error) {
	pt, err := m.db.GetPlaylistTrack(ctx, id)
	if err != nil {
		return types.PlaylistTrack{}, err
	}

	plist, err := m.db.GetPlaylist(ctx, pt.PlaylistID, userID)
	if err != nil {
		return types.PlaylistTrack{}, err
	}

	if plist.Owner != userID {
		return types.PlaylistTrack{}, errors.New("user is not the owner of the selected playlist")
	}

	return pt, nil
}
