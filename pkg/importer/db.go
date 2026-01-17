package importer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (i Importer) addAlbumTracks(ctx context.Context, tx database.DatabaseFace, albumTracks []types.AlbumTrack, albumID uuid.UUID) error {
	for _, at := range albumTracks {
		at.AlbumID = albumID
		exists, err := tx.AlbumTrackExists(ctx, at.AlbumID, at.TrackID)
		if err != nil {
			return err
		}
		if !exists {
			if _, err := tx.AddAlbumTrack(ctx, at); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i Importer) addAlbumArtists(ctx context.Context, tx database.DatabaseFace, albumID, artistID uuid.UUID) error {
	role := types.ArPrimary

	exists, err := tx.AlbumArtistExists(ctx, albumID, artistID, role)
	if err != nil {
		return err
	}

	aa := types.AlbumArtist{
		AlbumID:  albumID,
		ArtistID: artistID,
		Role:     role,
		Position: 0,
	}

	if !exists {
		if _, err = tx.AddAlbumArtist(ctx, aa); err != nil {
			return err
		}
	}

	return nil
}

func (i Importer) addOrGetArtist(ctx context.Context, tx database.DatabaseFace, artist types.Artist) (id uuid.UUID, retArtist *types.Artist, err error) {
	var existingArtist types.Artist

	if artist.MusicBrainzID != nil {
		existingArtist, err = tx.GetArtistFromMbID(ctx, *artist.MusicBrainzID)
	} else {
		err = sql.ErrNoRows
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingArtist, err = tx.GetArtistFromName(ctx, artist.Name)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return uuid.Nil, nil, err
				}
			} else {
				i.log.InfoContext(ctx, "found existsing artist using name", "id", existingArtist.ID)
				return existingArtist.ID, &existingArtist, nil
			}
		} else {
			return uuid.Nil, nil, err
		}
	} else {
		i.log.InfoContext(ctx, "found existsing artist using musicbrainz id", "id", existingArtist.ID)
		return existingArtist.ID, &existingArtist, nil
	}

	i.log.InfoContext(ctx, "creating new artist")
	id, err = tx.AddArtist(ctx, artist)

	return id, nil, err
}

func (i Importer) addOrUpdateTrack(ctx context.Context, tx database.DatabaseFace, track types.Track, albumID uuid.UUID) (id uuid.UUID, new bool, err error) {
	var existingTrack types.Track

	if track.MusicBrainzID != nil {
		existingTrack, err = tx.GetTrackFromMbID(ctx, *track.MusicBrainzID)
	} else {
		err = sql.ErrNoRows
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingTrack, err = tx.GetTrackFromName(ctx, albumID, track.Title)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return uuid.Nil, false, err
				}
			} else {
				i.log.InfoContext(ctx, "found existing track using name", "id", existingTrack.ID)
			}
		} else {
			return uuid.Nil, false, err
		}
	} else {
		i.log.InfoContext(ctx, "found existing track using musicbrainz id", "id", existingTrack.ID)
	}

	if existingTrack.ID != uuid.Nil {
		track.ID = existingTrack.ID
		err = tx.UpdateTrack(ctx, track)
		if err != nil {
			return uuid.Nil, false, err
		}
		return track.ID, false, nil
	} else {
		i.log.InfoContext(ctx, "creating new track")
		id, err := tx.AddTrack(ctx, track)
		return id, true, err
	}
}
