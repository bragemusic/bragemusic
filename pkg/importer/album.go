package importer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bragemusic/bragemusic/pkg/acoustid"
	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/musicbrainz"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/bragemusic/bragemusic/pkg/utils"
	"github.com/dhowden/tag"
	"github.com/gofrs/uuid/v5"
)

var (
	ErrAlbumMbIDNotFound  = errors.New("could not find MusicBrainz ID for album")
	ErrArtistMbIDNotFound = errors.New("could not find MusicBrainz ID for artist")
	ErrId3AlbumNotFound   = errors.New("could not get album name from ID3")
	ErrAlbumNotFound      = errors.New("no existing album found")
)

func (i Importer) addAlbum(ctx context.Context, tx database.DatabaseFace, albumAnalysis AlbumAnalysisResults, existingAlbum types.Album, userID uuid.UUID) (album types.Album, err error) {
	if albumAnalysis.AlbumID != "" {
		album, err = i.generateAlbumFromMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			return types.Album{}, err
		}
	} else {
		album = i.generateAlbumFromID3(ctx, albumAnalysis)
	}

	// Issue #140: Should fix albums with empty string as mbID
	if album.MusicBrainzID != nil && *album.MusicBrainzID == "" {
		album.MusicBrainzID = nil
	}

	if existingAlbum.ID != uuid.Nil {
		// Exisiting album is not created from musicbrainz id but new one is, update it
		if existingAlbum.MusicBrainzID == nil && album.MusicBrainzID != nil {
			album.ID = existingAlbum.ID
			if err = tx.UpdateAlbum(ctx, album, userID); err != nil {
				return types.Album{}, err
			}
		} else {
			album = existingAlbum
		}
	} else {
		album.ID, err = tx.AddAlbum(ctx, album, userID)
		if err != nil {
			return types.Album{}, err
		}
	}

	return album, nil
}

type MbAlbum struct {
	AlbumID           string
	MbTrackIDs        []*string
	MatchedTracks     int
	ID3AlbumNameMatch float32
	AIDMatch          acoustid.AcoustMatch
	ReleaseDate       acoustid.Date
}

type AlbumAnalysisResults struct {
	AlbumID           string
	Id3Artist         string
	Id3Album          string
	Id3ReleaseDate    int
	MatchedTracks     int
	ID3AlbumNameMatch float32
	MediaFiles        []types.MediaFile
	Tracks            []Track
	Covers            []*tag.Picture
}

type Track struct {
	TrackNumber *int
	DiscNumber  *int
	Name        *string
	MbID        *string
	MediaFileID uuid.UUID
}

func (i Importer) generateAlbumFromMbID(ctx context.Context, mbID string) (types.Album, error) {
	mbAlbum, err := i.mb.GetAlbum(ctx, mbID)
	if err != nil {
		return types.Album{}, err
	}

	album := types.Album{
		AlbumBase: types.AlbumBase{
			MusicBrainzID: &mbAlbum.ID,
			Name:          mbAlbum.Title,
			SortName:      mbAlbum.Title,
			ReleaseDate:   &mbAlbum.Date.Time,
			Tracks:        &mbAlbum.TrackCount,
		},
	}

	if len(mbAlbum.Media) > 0 {
		album.Discs = &mbAlbum.Media[0].DiscCount
	}

	if album.Discs != nil && *album.Discs == 0 {
		album.Discs = utils.Ptr(1)
	}

	return album, nil
}

func (i Importer) generateAlbumFromID3(ctx context.Context, analysResults AlbumAnalysisResults) types.Album {
	var releaseDate *time.Time

	rD, err := time.Parse("2006", fmt.Sprint(analysResults.Id3ReleaseDate))
	if err == nil {
		releaseDate = &rD
	}

	totalTracks := len(analysResults.Tracks)
	// FIXME do something
	totalDiscs := 1

	albumName := analysResults.Id3Album
	if strings.TrimSpace(albumName) == "" {
		albumName = "Unknown Album"
	}

	album := types.Album{
		AlbumBase: types.AlbumBase{
			Name:        albumName,
			SortName:    albumName,
			ReleaseDate: releaseDate,
			Tracks:      &totalTracks,
			Discs:       &totalDiscs,
		},
	}

	return album
}

func (i Importer) generateArtistFromAlbumMbID(ctx context.Context, mbID string) (types.Artist, error) {
	mbAlbum, err := i.mb.GetAlbum(ctx, mbID)
	if err != nil {
		return types.Artist{}, err
	}

	if len(mbAlbum.ArtistCredit) == 0 {
		return types.Artist{}, ErrArtistMbIDNotFound
	}

	artistMbID := mbAlbum.ArtistCredit[0].Artist.ID

	mbArtist, err := i.mb.GetArtist(ctx, artistMbID)
	if err != nil {
		return types.Artist{}, err
	}

	artist := types.Artist{
		ArtistBase: types.ArtistBase{
			MusicBrainzID: &mbArtist.ID,
			Name:          mbArtist.Name,
			SortName:      mbArtist.SortName,
		},
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

	return artist, nil
}

func (i Importer) generateArtistFromID3(ctx context.Context, analysis AlbumAnalysisResults) types.Artist {
	artist := types.Artist{
		ArtistBase: types.ArtistBase{
			Name:     analysis.Id3Artist,
			SortName: analysis.Id3Artist,
		},
	}

	return artist
}

func (i Importer) generateTracksFromAlbumMbID(ctx context.Context, albumMbID string) (tracks []types.Track, albumTracks []types.AlbumTrack, err error) {
	mbAlbum, err := i.mb.GetAlbum(ctx, albumMbID)
	if err != nil {
		return nil, nil, err
	}

	for discIdx, media := range mbAlbum.Media {
		for _, mbTrack := range media.Tracks {
			tr := i.generateTrack(mbTrack)
			atr := i.generateAlbumTrack(mbTrack, discIdx+1)
			tracks = append(tracks, tr)
			albumTracks = append(albumTracks, atr)
		}
	}

	return tracks, albumTracks, nil
}

func (i Importer) generateTrack(mbTrack musicbrainz.Track) types.Track {
	return types.Track{
		Title:         mbTrack.Title,
		MusicBrainzID: &mbTrack.ID,
	}
}

func (i Importer) generateAlbumTrack(mbTrack musicbrainz.Track, discNumber int) types.AlbumTrack {
	return types.AlbumTrack{
		DiscNumber:  discNumber,
		TrackNumber: mbTrack.Position,
	}
}

func (i Importer) getID3Info(ctx context.Context, files []types.MediaFile) (artist, album string, releaseYear int, tracks []Track, pics []*tag.Picture, err error) {
	albums := []string{}
	artists := []string{}
	releaseYears := []int{}

	for _, f := range files {
		filename := filepath.Join(i.musicDir, f.Filename())
		r, err := os.Open(filename)
		if err != nil {
			return "", "", 0, nil, nil, err
		}

		md, err := tag.ReadFrom(r)
		if err != nil {
			return "", "", 0, nil, nil, err
		}

		if md.Album() != "" {
			albums = append(albums, md.Album())
		}

		if md.Artist() != "" {
			artists = append(artists, md.Artist())
		}
		releaseYears = append(releaseYears, md.Year())

		track := Track{}
		track.MediaFileID = f.ID
		tn, _ := md.Track()
		dn, _ := md.Disc()

		discNumber, trackNumber, ok := utils.ExtractDiscAndTrack(f.OrgFilename)

		if tn != 0 {
			track.TrackNumber = &tn
		} else {
			if ok {
				track.TrackNumber = utils.Ptr(trackNumber)
			} else {
				track.TrackNumber = utils.Ptr(1)
			}
		}

		if dn != 0 {
			track.DiscNumber = &dn
		} else {
			if ok {
				track.DiscNumber = utils.Ptr(discNumber)
			} else {
				track.DiscNumber = utils.Ptr(1)
			}
		}
		track.Name = utils.Ptr(md.Title())

		tracks = append(tracks, track)
		pics = append(pics, md.Picture())

		r.Close()
	}

	i.log.DebugContext(ctx, "found album names in ID3", "names", albums)

	if len(albums) == 0 {
		return "", "", 0, nil, nil, ErrId3AlbumNotFound
	}

	return utils.HighestCount(artists), utils.HighestCount(albums), utils.HighestCount(releaseYears), tracks, pics, nil
}

func (i Importer) getExistingAlbum(ctx context.Context, albumAnalysis AlbumAnalysisResults) (types.Album, error) {
	if albumAnalysis.AlbumID != "" {
		album, err := i.db.GetAlbumFromMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return types.Album{}, err
			}
		} else {
			return album, nil
		}
	}

	album, err := i.db.GetAlbumFromArtistAndName(ctx, albumAnalysis.Id3Artist, albumAnalysis.Id3Album)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return types.Album{}, err
		}
	} else {
		return album, nil
	}

	return types.Album{}, ErrAlbumNotFound
}
