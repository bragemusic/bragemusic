package importer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
)

func (i Importer) addOrGetArtist(ctx context.Context, tx database.DatabaseFace, artist types.Artist) (id string, retArtist *types.Artist, err error) {
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
					return "", nil, err
				}
			} else {
				i.log.InfoContext(ctx, "found existsing artist using name", "id", existingArtist.ID)
				return existingArtist.ID, &existingArtist, nil
			}
		} else {
			return "", nil, err
		}
	} else {
		i.log.InfoContext(ctx, "found existsing artist using musicbrainz id", "id", existingArtist.ID)
		return existingArtist.ID, &existingArtist, nil
	}

	i.log.InfoContext(ctx, "creating new artist")
	id, err = tx.AddArtist(ctx, artist)

	return id, nil, err
}

func (i Importer) addOrGetAlbum(ctx context.Context, tx database.DatabaseFace, album types.Album, artistName string) (id string, err error) {
	var existingAlbum types.Album

	if album.MusicBrainzID != nil {
		existingAlbum, err = tx.GetAlbumFromMbID(ctx, *album.MusicBrainzID)
	} else {
		err = sql.ErrNoRows
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingAlbum, err = tx.GetAlbumFromArtistAndName(ctx, artistName, album.Name)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return "", err
				}
			} else {
				i.log.InfoContext(ctx, "found existsing album using name", "id", existingAlbum.ID)
				return existingAlbum.ID, nil
			}
		} else {
			return "", err
		}
	} else {
		i.log.InfoContext(ctx, "found existsing album using musicbrainz id", "id", existingAlbum.ID)
		return existingAlbum.ID, nil
	}

	i.log.InfoContext(ctx, "creating new album")
	return tx.AddAlbum(ctx, album)
}

func (i Importer) addOrUpdateTrack(ctx context.Context, tx database.DatabaseFace, track types.Track) (id string, err error) {
	var existingTrack types.Track

	if track.MusicBrainzID != nil {
		existingTrack, err = tx.GetTrackFromMbID(ctx, *track.MusicBrainzID)
	} else {
		err = sql.ErrNoRows
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingTrack, err = tx.GetTrackFromName(ctx, *track.AlbumID, track.Title)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return "", err
				}
			} else {
				i.log.InfoContext(ctx, "found existing track using name", "id", existingTrack.ID)
			}
		} else {
			return "", err
		}
	} else {
		i.log.InfoContext(ctx, "found existing track using musicbrainz id", "id", existingTrack.ID)
	}

	if existingTrack.ID != "" {
		track.ID = existingTrack.ID
		err = tx.UpdateTrack(ctx, track)
		if err != nil {
			return "", err
		}
		return track.ID, nil
	} else {
		i.log.InfoContext(ctx, "creating new track")
		return tx.AddTrack(ctx, track)
	}
}
