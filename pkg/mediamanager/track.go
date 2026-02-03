package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (m MediaManager) GetTrack(ctx context.Context, trackID uuid.UUID) (types.Track, error) {
	track, err := m.db.GetTrackFromID(ctx, trackID)
	if err != nil {
		return track, m.berr.DatabaseError(err, types.EntityTrack, &trackID)
	}

	return track, nil
}

func (m MediaManager) ListTracksByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.Track, error) {
	tracks, err := m.db.GetTracksFromAlbumID(ctx, albumID)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) ListTracksDetailedByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.TrackDetailed, error) {
	tracks, err := m.db.ListAlbumTracksDetailed(ctx, albumID)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) ListTracksDetailedByArtist(ctx context.Context, artistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder, limit *int, includeMissingFiles bool) ([]types.TrackDetailed, error) {
	tracks, err := m.db.GetTracksDetailedFromArtistID(ctx, artistID, sortBy, sortOrder, limit, includeMissingFiles)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) UpdateTrack(ctx context.Context, trackID uuid.UUID, trackData types.TrackUpdate) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityTrack, &trackID)
	}
	defer tx.Rollback()

	existingTrack, err := tx.GetTrackFromID(ctx, trackID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityTrack, &trackID)
	}

	trackData.ID = trackID
	trackData.MediaFile = existingTrack.MediaFile

	err = tx.UpdateTrack(ctx, trackData.Track)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityTrack, &trackID)
	}

	albumTrack, err := tx.GetAlbumTrackFromAlbumAndTrack(ctx, trackData.AlbumID, trackData.ID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityAlbumTrack, nil)
	}

	if albumTrack.DiscNumber != trackData.DiscNumber || albumTrack.TrackNumber != trackData.TrackNumber {
		albumTrack.DiscNumber = trackData.DiscNumber
		albumTrack.TrackNumber = trackData.TrackNumber

		if err = tx.UpdateAlbumTrack(ctx, albumTrack); err != nil {
			return m.berr.DatabaseError(err, types.EntityAlbumTrack, &albumTrack.ID)
		}
	}

	return tx.Commit()
}

func (m MediaManager) CountTracks(ctx context.Context) (int, error) {
	cnt, err := m.db.CountTracks(ctx)
	if err != nil {
		return 0, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return cnt, nil
}
