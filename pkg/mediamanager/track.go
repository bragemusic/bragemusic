package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) GetTrack(ctx context.Context, trackID string) (types.Track, error) {
	track, err := m.db.GetTrackFromID(ctx, trackID)
	if err != nil {
		return track, err
	}

	return track, nil
}

func (m MediaManager) ListTracksByAlbum(ctx context.Context, albumID string) ([]types.Track, error) {
	tracks, err := m.db.GetTracksFromAlbumID(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (m MediaManager) ListTracksDetailedByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.TrackDetailed, error) {
	tracks, err := m.db.ListAlbumTracksDetailed(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (m MediaManager) ListTracksDetailedByArtist(ctx context.Context, artistID string, sortBy database.SortBy, sortOrder database.SortOrder, limit *int, includeMissingFiles bool) ([]types.TrackDetailed, error) {
	tracks, err := m.db.GetTracksDetailedFromArtistID(ctx, artistID, sortBy, sortOrder, limit, includeMissingFiles)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (m MediaManager) UpdateTrack(ctx context.Context, trackID uuid.UUID, trackData types.TrackUpdate) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existingTrack, err := tx.GetTrackFromID(ctx, trackID.String())
	if err != nil {
		return err
	}

	trackData.ID = trackID
	trackData.MediaFile = existingTrack.MediaFile

	err = tx.UpdateTrack(ctx, trackData.Track)
	if err != nil {
		return err
	}

	albumTrack, err := tx.GetAlbumTrackFromAlbumAndTrack(ctx, trackData.AlbumID, trackData.ID)
	if err != nil {
		return err
	}

	if albumTrack.DiscNumber != trackData.DiscNumber || albumTrack.TrackNumber != trackData.TrackNumber {
		albumTrack.DiscNumber = trackData.DiscNumber
		albumTrack.TrackNumber = trackData.TrackNumber

		if err = tx.UpdateAlbumTrack(ctx, albumTrack); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m MediaManager) CountTracks(ctx context.Context) (int, error) {
	return m.db.CountTracks(ctx)
}
