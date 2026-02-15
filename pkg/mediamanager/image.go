package mediamanager

import (
	"context"
	"path/filepath"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddArtistImage(ctx context.Context, filename string, artistID, userID uuid.UUID) error {
	if m.im == nil {
		return m.berr.DependencyMissing(nil, "imagemagick")
	}

	artist, err := m.GetArtist(ctx, artistID)
	if err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "artists", artistID.String())

	if err = m.im.ResizeAll(ctx, filename, outputFolder); err != nil {
		return err
	}

	if err = m.db.UpdateArtist(ctx, artist, userID); err != nil {
		return m.berr.DatabaseError(err, types.EntityArtist, &artistID)
	}

	return nil
}

func (m MediaManager) AddAlbumImage(ctx context.Context, filename string, albumID uuid.UUID) error {
	if m.im == nil {
		return m.berr.DependencyMissing(nil, "imagemagick")
	}

	album, err := m.GetAlbum(ctx, albumID)
	if err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "albums", albumID.String())

	if err = m.im.ResizeAll(ctx, filename, outputFolder); err != nil {
		return err
	}

	if err = m.db.UpdateAlbum(ctx, album); err != nil {
		return m.berr.DatabaseError(err, types.EntityAlbum, &albumID)
	}

	return nil
}

func (m MediaManager) AddPlaylistImage(ctx context.Context, filename string, id uuid.UUID) error {
	if m.im == nil {
		return m.berr.DependencyMissing(nil, "imagemagick")
	}

	plist, err := m.GetPlaylist(ctx, id)
	if err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "playlists", id.String())

	if err = m.im.ResizeAll(ctx, filename, outputFolder); err != nil {
		return err
	}

	if err = m.db.UpdatePlaylist(ctx, plist); err != nil {
		return m.berr.DatabaseError(err, types.EntityPlaylist, &id)
	}

	return nil
}
