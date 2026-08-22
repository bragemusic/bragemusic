package mediamanager

import (
	"context"

	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

func (m MediaManager) GetTrack(ctx context.Context, trackID uuid.UUID) (types.Track, error) {
	track, err := m.db.GetTrackFromID(ctx, trackID)
	if err != nil {
		return track, m.berr.DatabaseError(err, types.EntityTrack, &trackID)
	}

	return track, nil
}

func (m MediaManager) GetTrackAnalysisByID(ctx context.Context, id uuid.UUID) (types.TrackAnalysis, error) {
	trackAnalysis, err := m.db.GetTrackAnalysisByID(ctx, id)
	if err != nil {
		return types.TrackAnalysis{}, m.berr.DatabaseError(err, types.EntityTrackAnalysis, &id)
	}

	return trackAnalysis, nil
}

func (m MediaManager) GetTrackArtistByID(ctx context.Context, id uuid.UUID) (types.TrackArtist, error) {
	trackArtist, err := m.db.GetTrackArtistByID(ctx, id)
	if err != nil {
		return types.TrackArtist{}, m.berr.DatabaseError(err, types.EntityTrackArtist, &id)
	}

	return trackArtist, nil
}

func (m MediaManager) GetTrackDetailed(ctx context.Context, trackID, albumID, userID uuid.UUID) (types.TrackDetailed, error) {
	track, err := m.db.GetTrackDetailed(ctx, trackID, albumID, userID)
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

func (m MediaManager) ListTracksDetailedByAlbum(ctx context.Context, albumID, userID uuid.UUID) ([]types.TrackDetailed, error) {
	tracks, err := m.db.ListAlbumTracksDetailed(ctx, albumID, userID)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) ListTracksDetailedByArtist(ctx context.Context, artistID, userID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder, limit *int, includeMissingFiles bool) ([]types.TrackDetailed, error) {
	tracks, err := m.db.GetTracksDetailedFromArtistID(ctx, artistID, userID, sortBy, sortOrder, limit, includeMissingFiles)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	return tracks, nil
}

func (m MediaManager) ListTracksDetailed(ctx context.Context, filter types.TrackFilter, page, limit int) (items []types.TrackDetailed, actualPage, actualLimit, totalPages, totalItems int, err error) {
	page = max(1, page)
	limit = min(max(10, limit), 10000000)

	tracks, totalPages, totalItems, err := m.db.ListTracksWithFilters(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, 0, 0, 0, m.berr.DatabaseError(err, types.EntityTrack, nil)
	}

	for _, t := range tracks {
		items = append(items, types.TrackDetailed{
			ID:        t.Track.ID,
			Title:     t.Track.Title,
			AlbumID:   t.Album.ID.String(),
			AlbumName: t.Album.Name,
			Artists: lo.Map(t.Artists, func(item types.Artist, index int) types.ArtistMinimal {
				return types.ArtistMinimal{
					ID:       item.ID,
					Name:     item.Name,
					SortName: item.SortName,
					Role:     item.Role,
				}
			}),
			// ArtistIDs:     lo.Map(t.Artists, func(item types.Artist, index int) string { return item.ID.String() }),
			// ArtistNames:   lo.Map(t.Artists, func(item types.Artist, index int) string { return item.Name }),
			MusicBrainzID: t.Track.MusicBrainzID,
			TrackNumber:   t.AlbumTrack.TrackNumber,
			DiscNumber:    t.AlbumTrack.DiscNumber,
			Genre:         nil,
			Comment:       nil,
			MediaFile:     t.Mediafile,
			PlayCount:     0,
			ContextID:     nil,
			Rating:        nil,
			UserRating:    nil,
			Liked:         false,
			CreatedAt:     t.Track.CreatedAt,
			UpdatedAt:     t.Track.UpdatedAt,
			Analysis:      t.Analysis,
		})
	}

	return items, page, limit, totalPages, totalItems, nil
}

func (m MediaManager) UpdateTrack(ctx context.Context, trackID uuid.UUID, trackData types.TrackUpdate, userID uuid.UUID) error {
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

	err = tx.UpdateTrack(ctx, trackData.Track, userID)
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

		if err = tx.UpdateAlbumTrack(ctx, albumTrack, userID); err != nil {
			return m.berr.DatabaseError(err, types.EntityAlbumTrack, &albumTrack.ID)
		}
	}

	// Track artists
	existingTrackArtists, err := tx.ListTrackArtistsByTrackID(ctx, trackID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityTrackArtist, nil)
	}

	for _, artist := range trackData.Artists {
		exists := lo.ContainsBy(existingTrackArtists, func(item types.TrackArtist) bool {
			return item.ArtistID == artist
		})

		if !exists {
			ta := types.TrackArtist{
				TrackID:  trackID,
				ArtistID: artist,
				Role:     types.ArFeatured,
			}
			if _, err := tx.AddTrackArtist(ctx, ta, userID); err != nil {
				return m.berr.DatabaseError(err, types.EntityTrackArtist, nil)
			}
		}
	}

	for _, existingTa := range existingTrackArtists {
		exists := lo.ContainsBy(trackData.Artists, func(item uuid.UUID) bool {
			return item == existingTa.ArtistID
		})

		if !exists {
			if err := tx.DeleteTrackArtist(ctx, existingTa.ID, userID); err != nil {
				return m.berr.DatabaseError(err, types.EntityTrackArtist, &existingTa.ID)
			}
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
