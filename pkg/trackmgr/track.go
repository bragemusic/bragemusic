package trackmgr

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
)

func (t TrackManager) generateTrack(mbTrack musicbrainz.Track, albumID string, discNumber int) types.Track {
	return types.Track{
		Title:         mbTrack.Title,
		AlbumID:       &albumID,
		MusicBrainzID: &mbTrack.ID,
		TrackNumber:   &mbTrack.Position,
		DiscNumber:    &discNumber,
	}
}

func (t TrackManager) generateTrackFromID3(metadata tag.Metadata, albumID string) types.Track {
	trackNumber, _ := metadata.Track()
	discNumber, _ := metadata.Disc()

	track := types.Track{
		Title:       metadata.Title(),
		AlbumID:     &albumID,
		TrackNumber: &trackNumber,
		DiscNumber:  &discNumber,
		Genre:       utils.Ptr(metadata.Genre()),
		Year:        utils.Ptr(metadata.Year()),
		Composer:    utils.Ptr(metadata.Composer()),
		Comment:     utils.Ptr(metadata.Comment()),
	}

	if metadata.Artist() != metadata.AlbumArtist() {
		track.TrackArtist = utils.Ptr(metadata.Artist())
	}

	return track
}

func (t TrackManager) generateTracks(ctx context.Context, tx database.DatabaseFace, album types.Album, metadata tag.Metadata) (tracks []types.Track, new bool, err error) {
	// we have a MusicBrainz ID, so we do all if it with it
	if album.MusicBrainzID != nil {
		mbAlbum, err := t.mb.GetAlbum(ctx, *album.MusicBrainzID)
		if err != nil {
			return nil, false, err
		}

		// If we find the tracks for the album ID, and the album has a MusicBrainz ID we can assume that all tracks are there
		tracks, err = tx.GetTracksFromAlbumID(ctx, album.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, false, err
			}
		} else if len(tracks) > 0 {
			t.log.InfoContext(ctx, "tracks found in db using MusicBrainz ID", "album", album.Name)
			return tracks, false, nil
		}

		// Otherwise we generate the tracks from MusicBrainz data
		for discIdx, media := range mbAlbum.Media {
			for _, mbTrack := range media.Tracks {
				tr := t.generateTrack(mbTrack, album.ID, discIdx+1)
				tracks = append(tracks, tr)
			}
		}

		t.log.InfoContext(ctx, "tracks generated using MusicBrainz ID", "album", album.Name)
		return tracks, true, nil
	}

	// See if we have a track with the same name on the same album
	track, err := tx.GetTrackFromName(ctx, album.ID, metadata.Title())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, false, err
		}
	} else if len(tracks) > 0 {
		t.log.InfoContext(ctx, "track found in db using ID3 name", "album", album.Name)
		return []types.Track{track}, false, nil
	}

	track = t.generateTrackFromID3(metadata, album.ID)
	t.log.InfoContext(ctx, "track generated using ID3", "album", album.Name)

	return []types.Track{track}, true, nil
}

func (t TrackManager) updateTrackData(track types.Track, audioFile types.AudioFile, filename string, filetype tag.FileType) types.Track {
	track.DurationMS = utils.Ptr(audioFile.DurationMS())
	track.Bitrate = utils.Ptr(audioFile.Bitrate())
	track.SampleRate = utils.Ptr(audioFile.SampleRate())
	track.FilePath = filename
	track.FileSize = utils.Ptr(audioFile.FileSize())
	track.MimeType = utils.Ptr(string(filetype))

	return track
}

func (t TrackManager) GetTracksByAlbum(ctx context.Context, albumID string) ([]types.Track, error) {
	tracks, err := t.db.GetTracksFromAlbumID(ctx, albumID)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}
