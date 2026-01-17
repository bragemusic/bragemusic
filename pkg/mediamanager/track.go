package mediamanager

import (
	"context"
	"fmt"

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

	fmt.Println(trackData.DiscNumber, trackData.TrackNumber)

	// title = :title,
	// musicbrainz_id = :musicbrainz_id,
	// genre = :genre,
	// comment = :comment,
	// media_file = :media_file,
	// updated_at = :updated_at
	err = tx.UpdateTrack(ctx, trackData.Track)
	if err != nil {
		return err
	}

	// existingAlbumArtists, err := tx.ListAlbumArtistsByAlbumID(ctx, albumID)
	// if err != nil {
	// 	return err
	// }

	// for _, artistID := range albumData.Artists {
	// 	exists := lo.ContainsBy(existingAlbumArtists, func(item types.AlbumArtist) bool {
	// 		return item.ArtistID == artistID
	// 	})

	// 	if !exists {
	// 		aa := types.AlbumArtist{
	// 			AlbumID:  albumID,
	// 			ArtistID: artistID,
	// 			Role:     types.ArPrimary,
	// 			// FIXME: Need to do something about positions
	// 			Position: 0,
	// 		}
	// 		if _, err := tx.AddAlbumArtist(ctx, aa); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// for _, existingAa := range existingAlbumArtists {
	// 	exists := lo.ContainsBy(albumData.Artists, func(item uuid.UUID) bool {
	// 		return item == existingAa.ArtistID
	// 	})

	// 	if !exists {
	// 		if err := tx.DeleteAlbumArtist(ctx, existingAa.ID); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return tx.Commit()
}
