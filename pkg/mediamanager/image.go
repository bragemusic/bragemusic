package mediamanager

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) AddArtistImage(ctx context.Context, filename string, artistID uuid.UUID) error {
	if m.im == nil {
		return errors.New("imagemagick is not loaded")
	}

	artist, err := m.db.GetArtistFromID(ctx, artistID.String())
	if err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "artists", artistID.String())

	if err = m.im.ResizeAll(ctx, filename, outputFolder); err != nil {
		return err
	}

	if err = m.db.UpdateArtist(ctx, artist); err != nil {
		return err
	}

	return nil
}

func (m MediaManager) AddAlbumImage(ctx context.Context, filename string, albumID uuid.UUID) error {
	if m.im == nil {
		return errors.New("imagemagick is not loaded")
	}

	album, err := m.db.GetAlbumFromID(ctx, albumID.String())
	if err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "albums", albumID.String())

	if err = m.im.ResizeAll(ctx, filename, outputFolder); err != nil {
		return err
	}

	if err = m.db.UpdateAlbum(ctx, album); err != nil {
		return err
	}

	return nil
}

func (m MediaManager) AddPlaylistImage(ctx context.Context, filename string, id uuid.UUID) error {
	if m.im == nil {
		return errors.New("imagemagick is not loaded")
	}

	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return err
	}

	plist, err := m.db.GetPlaylist(ctx, id, user.ID)
	if err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "playlists", id.String())

	if err = m.im.ResizeAll(ctx, filename, outputFolder); err != nil {
		return err
	}

	if err = m.db.UpdatePlaylist(ctx, plist); err != nil {
		return err
	}

	return nil
}
