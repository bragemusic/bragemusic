package mediamanager

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) GetTrack(ctx context.Context, trackID string) (types.Track, error) {
	track, err := m.db.GetTrackFromID(ctx, trackID)
	if err != nil {
		return track, err
	}

	return track, nil
}

func (m MediaManager) GetTrackFile(ctx context.Context, trackID string, w io.Writer) error {
	track, err := m.db.GetTrackFromID(ctx, trackID)
	if err != nil {
		return err
	}

	if track.FilePath == "" {
		return fmt.Errorf("file does not exist for track '%s'", trackID)
	}

	f, err := os.Open(filepath.Join(m.musicDir, track.FilePath))
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

func (m MediaManager) ListTracksByAlbum(ctx context.Context, albumID string) ([]types.Track, error) {
	tracks, err := m.db.GetTracksFromAlbumID(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (m MediaManager) ListEnhancedTracksByAlbum(ctx context.Context, albumID string) ([]types.TrackEnhanced, error) {
	tracks, err := m.db.GetEnhancedTracksFromAlbumID(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}
