package trackmgr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/dhowden/tag"
)

func (t TrackManager) generateAlbumFromId3(metadata tag.Metadata, artistID string) types.Album {
	var releaseDate *time.Time

	rD, err := time.Parse("2006", fmt.Sprint(metadata.Year()))
	if err == nil {
		releaseDate = &rD
	}

	_, totalTracks := metadata.Track()
	_, totalDiscs := metadata.Disc()

	albumName := metadata.Album()
	if strings.TrimSpace(albumName) == "" {
		albumName = "Unknown Album"
	}

	album := types.Album{
		Name:        albumName,
		SortName:    albumName,
		ArtistID:    artistID,
		ReleaseDate: releaseDate,
		Tracks:      &totalTracks,
		Discs:       &totalDiscs,
	}

	return album
}

func (t TrackManager) generateAlbum(mbAlbum musicbrainz.Release, artistID string) types.Album {
	album := types.Album{
		MusicBrainzID: &mbAlbum.ID,
		Name:          mbAlbum.Title,
		SortName:      mbAlbum.Title,
		ArtistID:      artistID,
		ReleaseDate:   &mbAlbum.Date.Time,
		Tracks:        &mbAlbum.TrackCount,
	}

	if len(mbAlbum.Media) > 0 {
		album.Discs = &mbAlbum.Media[0].DiscCount
	}

	// TODO: Add description to artist from somewhere
	// TODO: add owner
	// TODO: add public

	return album
}

func (t TrackManager) getOrCreateAlbum(ctx context.Context, tx database.DatabaseFace, aIdMatches []acoustid.AcoustMatch, metadata tag.Metadata) (album types.Album, new bool, err error) {
	// 	ta forst fram ett album. 1 kolla db, 2 skapa med aid, 3 skapa med id3
	//     ta fram artist likadant
	//     ta fram track likadant
	//     Fyll pa med musicbrains om det har dykt upp battre info an tidigare
	// Da kan jag gora lite battre funktion. 1 album funk, 1 artistfunk, 1 track funk osv

	albumsMbIds := []string{}
	for _, m := range aIdMatches {
		albumsMbIds = append(albumsMbIds, m.AlbumID)
	}

	if len(albumsMbIds) > 0 {
		// See if any of the matched albums exists in db
		albums, err := tx.GetAlbumsByMbIDs(ctx, albumsMbIds)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return types.Album{}, false, err
			}
		} else {
			if len(albums) > 0 {
				if len(albums) > 1 {
					t.log.WarnContext(ctx, "multiple albums found with the provided MB IDs. This will merge later", "name", albums[0].Name)
				}
				t.log.InfoContext(ctx, "album found in db using MusicBrainz ID", "name", albums[0].Name)
				return albums[0], false, nil
			}
		}
	}

	// Get artist from id3
	id3ArtistName := metadata.Artist()
	if metadata.AlbumArtist() != "" && metadata.AlbumArtist() != id3ArtistName {
		id3ArtistName = metadata.AlbumArtist()
	}

	// See if the album exists on a different or no musicbrainz id
	namedAlbum, err := tx.GetAlbumFromArtistAndName(ctx, id3ArtistName, metadata.Album())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return types.Album{}, false, nil
		}
	} else {
		t.log.InfoContext(ctx, "album found using album and artist names", "name", namedAlbum.Name)
		return namedAlbum, false, nil
	}

	// If we have any AcousticID matches, use them to create the album
	if len(aIdMatches) > 0 {
		aIdMatch := aIdMatches[0]
		filteredMatches := t.filterAcoustIdMatches(aIdMatches, id3ArtistName, metadata.Album())
		if len(filteredMatches) > 0 {
			aIdMatch = filteredMatches[0]
		}

		mbAlbum, err := t.mb.GetAlbum(ctx, aIdMatch.AlbumID)
		if err != nil {
			return types.Album{}, false, err
		}

		album := t.generateAlbum(mbAlbum, "")

		// Double check so an album with the same ID does not exist
		existingAlbum, err := tx.GetAlbumFromMbID(ctx, *album.MusicBrainzID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return types.Album{}, false, err
			}
		} else {
			t.log.InfoContext(ctx, "album found in db using MusicBrainz ID", "name", existingAlbum.Name)
			return existingAlbum, false, nil
		}

		existingAlbum, err = tx.GetAlbumFromArtistAndName(ctx, id3ArtistName, album.Name)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return types.Album{}, false, err
			}
		} else {
			t.log.InfoContext(ctx, "album found in db names", "name", existingAlbum.Name)
			return existingAlbum, false, nil
		}

		t.log.InfoContext(ctx, "album generated using AcousicID", "name", album.Name)
		return album, true, nil
	}

	// Try to get a musicbrainz match using ID3
	mbAlbum, err := t.mb.GetAlbumFromNames(ctx, id3ArtistName, metadata.Album())
	if err != nil {
		return types.Album{}, false, err
	}

	if mbAlbum != nil {
		album := t.generateAlbum(*mbAlbum, "")
		t.log.InfoContext(ctx, "album generated using MusicBrainz", "name", album.Name)
		return album, true, nil
	}

	// Finally if nothing else works, use id3
	album = t.generateAlbumFromId3(metadata, "")

	// Make sure it does not exist
	existingAlbum, err := tx.GetAlbumFromArtistAndName(ctx, id3ArtistName, album.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return types.Album{}, false, err
		}
	} else {
		t.log.InfoContext(ctx, "album found in db names", "name", existingAlbum.Name)
		return existingAlbum, false, nil
	}
	t.log.InfoContext(ctx, "album generated using ID3", "name", album.Name)

	return album, true, nil
}

func (t TrackManager) GetAlbumsByArtist(ctx context.Context, artistID string) ([]types.Album, error) {
	// albums, err := t.db.ListAlbumsByArtist(ctx, artistID)
	// if err != nil {
	// 	return nil, err
	// }

	// return albums, nil
	return nil, nil
}
