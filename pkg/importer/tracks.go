package importer

import (
	"context"
	"slices"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
)

func (i Importer) addMultipleTracks(ctx context.Context, tx database.DatabaseFace, albumAnalysis AlbumAnalysisResults, existingAlbumID uuid.UUID) ([]types.AlbumTrack, error) {
	var tracks []types.Track
	var albumTracks []types.AlbumTrack
	var err error

	if albumAnalysis.AlbumID != "" {
		tracks, albumTracks, err = i.generateTracksFromAlbumMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			return nil, err
		}
	}

	addedTrackIdx := []int{}

	for _, track := range albumAnalysis.Tracks {
		if track.MbID != nil {
			for tidx := range tracks {
				if *tracks[tidx].MusicBrainzID == *track.MbID {
					tracks[tidx].MediaFile = utils.Ptr(track.MediaFileID)

					tracks[tidx].ID, _, err = i.addOrUpdateTrack(ctx, tx, tracks[tidx], existingAlbumID)
					if err != nil {
						return nil, err
					}

					albumTracks[tidx].TrackID = tracks[tidx].ID
					addedTrackIdx = append(addedTrackIdx, tidx)
					break
				}
			}
		} else {
			t := types.Track{
				Title:     *track.Name,
				MediaFile: utils.Ptr(track.MediaFileID),
			}

			t.ID, _, err = i.addOrUpdateTrack(ctx, tx, t, existingAlbumID)
			if err != nil {
				return nil, err
			}

			tracks = append(tracks, t)
			albumTracks = append(albumTracks, types.AlbumTrack{
				DiscNumber:  *track.DiscNumber,
				TrackNumber: *track.TrackNumber,
				TrackID:     t.ID,
			})

			addedTrackIdx = append(addedTrackIdx, len(tracks)-1)
		}
	}

	// Add unmatched tracks
	for tidx := range tracks {
		if slices.Contains(addedTrackIdx, tidx) {
			continue
		}

		tracks[tidx].ID, _, err = i.addOrUpdateTrack(ctx, tx, tracks[tidx], existingAlbumID)
		if err != nil {
			return nil, err
		}

		albumTracks[tidx].TrackID = tracks[tidx].ID
	}

	return albumTracks, nil
}
