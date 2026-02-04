package mediamanager

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) GetMediaFile(ctx context.Context, mediafileID uuid.UUID) (types.MediaFile, error) {
	mf, err := m.db.GetMediaFile(ctx, mediafileID)
	if err != nil {
		return types.MediaFile{}, m.berr.DatabaseError(err, types.EntityMediaFile, &mediafileID)
	}

	return mf, nil
}

func (m MediaManager) GetMediaFileFile(ctx context.Context, mediafileID uuid.UUID, w io.Writer) error {
	mf, err := m.GetMediaFile(ctx, mediafileID)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Join(m.musicDir, mf.Filename()))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}
