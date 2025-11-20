package importer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
)

func (i Importer) addOrGetArtist(ctx context.Context, tx database.DatabaseFace, artist types.Artist) (id string, err error) {
	existingArtist, err := tx.GetArtistFromMbID(ctx, *artist.MusicBrainzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingArtist, err = tx.GetArtistFromName(ctx, artist.Name)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return "", err
				}
			} else {
				i.log.InfoContext(ctx, "found existsing artist using name", "id", existingArtist.ID)
				return existingArtist.ID, nil
			}
		} else {
			return "", err
		}
	} else {
		i.log.InfoContext(ctx, "found existsing artist using musicbrainz id", "id", existingArtist.ID)
		return existingArtist.ID, nil
	}

	i.log.InfoContext(ctx, "creating new artist")
	return tx.AddArtist(ctx, artist)
}

func (i Importer) addOrGetAlbum(ctx context.Context, tx database.DatabaseFace, album types.Album, artistName string) (id string, err error) {
	existingAlbum, err := tx.GetAlbumFromMbID(ctx, *album.MusicBrainzID)
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
