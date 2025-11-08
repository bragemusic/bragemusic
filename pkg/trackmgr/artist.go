package trackmgr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/bragemusic/core/pkg/wiki"
	"github.com/dhowden/tag"
)

func (t TrackManager) generateArtistFromId3(metadata tag.Metadata) types.Artist {
	artistName := metadata.Artist()
	if metadata.AlbumArtist() != metadata.Artist() {
		artistName = metadata.AlbumArtist()
	}

	artist := types.Artist{
		Name:     artistName,
		SortName: artistName,
	}

	return artist
}

func (t TrackManager) generateArtist(mbArtist musicbrainz.ArtistResponse) types.Artist {
	artist := types.Artist{
		MusicBrainzID: &mbArtist.ID,
		Name:          mbArtist.Name,
		SortName:      mbArtist.SortName,
	}

	if mbArtist.Country != "" {
		artist.Country = &mbArtist.Country
	}

	if mbArtist.LifeSpan.Begin != "" {
		begin, err := time.Parse("2006", mbArtist.LifeSpan.Begin)
		if err == nil {
			artist.YearStarted = utils.Ptr(begin.Year())
		}
	}

	if mbArtist.LifeSpan.Ended && mbArtist.LifeSpan.End != "" {
		end, err := time.Parse("2006", mbArtist.LifeSpan.End)
		if err == nil {
			artist.YearEnded = utils.Ptr(end.Year())
		}
	}

	// TODO: Add description to artist from somewhere

	return artist
}

func (t TrackManager) addOrGetArtistFromMbId(ctx context.Context, tx database.DatabaseFace, artistMbId string) (types.Artist, error) {
	artist, err := tx.GetArtistFromMbID(ctx, artistMbId)
	if err == nil {
		return artist, nil
	}

	mbArtist, err := t.mb.GetArtist(artistMbId)
	if err != nil {
		return types.Artist{}, err
	}

	for _, rel := range mbArtist.Relations {
		if rel.Type == "wikidata" {
			fmt.Println(rel)
		}
	}

	artist = t.generateArtist(mbArtist)
	artistID, err := tx.AddArtist(ctx, artist)
	if err != nil {
		return types.Artist{}, err
	}
	artist.ID = artistID

	return artist, nil
}

func (t TrackManager) GetArtistMetaData(ctx context.Context, artistMbId string) (wiki.WikiData, error) {
	mbArtist, err := t.mb.GetArtist(artistMbId)
	if err != nil {
		return wiki.WikiData{}, err
	}

	wikiDataUrl := ""

	for _, rel := range mbArtist.Relations {
		if rel.Type == "wikidata" {
			wikiDataUrl = rel.URL.Resource
		}
	}

	if wikiDataUrl == "" {
		return wiki.WikiData{}, fmt.Errorf("could not get wikidata for artist MbID with '%s'", artistMbId)
	}

	wikiData, err := t.wiki.GetWikiData(ctx, wikiDataUrl)
	if err != nil {
		return wiki.WikiData{}, err
	}

	return wikiData, nil
}

func (t TrackManager) GetAlbumMetaData(ctx context.Context, albumMbId string) (wiki.WikiData, error) {
	mbAlbum, err := t.mb.GetAlbum(albumMbId)
	if err != nil {
		return wiki.WikiData{}, err
	}

	wikiDataUrl := ""

	for _, rel := range mbAlbum.Relations {
		if rel.Type == "wikidata" {
			wikiDataUrl = rel.URL.Resource
		}
	}

	if wikiDataUrl == "" {
		return wiki.WikiData{}, fmt.Errorf("could not get wikidata for album MbID with '%s'", albumMbId)
	}

	wikiData, err := t.wiki.GetWikiData(ctx, wikiDataUrl)
	if err != nil {
		return wiki.WikiData{}, err
	}

	return wikiData, nil
}

func (t TrackManager) getOrCreateArtist(ctx context.Context, tx database.DatabaseFace, album types.Album, metadata tag.Metadata) (artist types.Artist, new bool, err error) {
	// 	ta forst fram ett album. 1 kolla db, 2 skapa med aid, 3 skapa med id3
	//     ta fram artist likadant
	//     ta fram track likadant
	//     Fyll pa med musicbrains om det har dykt upp battre info an tidigare
	// Da kan jag gora lite battre funktion. 1 album funk, 1 artistfunk, 1 track funk osv

	// Get artist from id3
	id3ArtistName := metadata.Artist()
	if metadata.AlbumArtist() != "" && metadata.AlbumArtist() != id3ArtistName {
		id3ArtistName = metadata.AlbumArtist()
	}

	// If MusicBrainz ID exists on the album, use that to get the artist
	if album.MusicBrainzID != nil {
		mbAlbum, err := t.mb.GetAlbum(*album.MusicBrainzID)
		if err != nil {
			return types.Artist{}, false, err
		}

		if len(mbAlbum.ArtistCredit) > 0 {
			dbArtist, err := tx.GetArtistFromMbID(ctx, mbAlbum.ArtistCredit[0].Artist.ID)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return types.Artist{}, false, err
				}
			} else {
				t.log.InfoContext(ctx, "artist found in db using MusicBrainz ID", "name", dbArtist.Name)
				return dbArtist, false, nil
			}
		}
	}

	// See if artist with the same name exists in db
	namedArtist, err := tx.GetArtistFromName(ctx, id3ArtistName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return types.Artist{}, false, err
		}
	} else {
		t.log.InfoContext(ctx, "artist found in db using ID3", "name", namedArtist.Name)
		return namedArtist, false, nil
	}

	// If MusicBrainzID exists, use that to create the artist
	if album.MusicBrainzID != nil {
		mbAlbum, err := t.mb.GetAlbum(*album.MusicBrainzID)
		if err != nil {
			return types.Artist{}, false, err
		}

		if len(mbAlbum.ArtistCredit) > 0 {
			mbArtist, err := t.mb.GetArtist(mbAlbum.ArtistCredit[0].Artist.ID)
			if err != nil {
				return types.Artist{}, false, err
			}

			artist := t.generateArtist(mbArtist)

			t.log.InfoContext(ctx, "artist generated using MusicBrainz ID", "name", artist.Name)
			return artist, true, nil
		}
	}

	// As a last resort, use ID3
	artist = t.generateArtistFromId3(metadata)
	t.log.InfoContext(ctx, "artist generated using ID3", "name", artist.Name)

	return artist, true, nil
}

func (t TrackManager) ListArtists(ctx context.Context) ([]types.Artist, error) {
	artists, err := t.db.ListArtists(ctx)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
