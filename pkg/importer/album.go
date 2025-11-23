package importer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/files"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
)

var (
	ErrAlbumMbIDNotFound  = errors.New("could not find MusicBrainz ID for album")
	ErrArtistMbIDNotFound = errors.New("could not find MusicBrainz ID for artist")
	ErrId3AlbumNotFound   = errors.New("could not get album name from ID3")
)

func (i Importer) importAlbumFiles(ctx context.Context, folder string) error {
	osFiles, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	filenames := []string{}
	for _, f := range osFiles {
		filenames = append(filenames, filepath.Join(folder, f.Name()))
	}

	var album types.Album

	albumAnalysis, err := i.analyzeAlbum(ctx, filenames)
	if err != nil {
		if errors.Is(err, ErrAlbumMbIDNotFound) {
			i.log.WarnContext(ctx, "no album musicbrainz ID found, using ID3")
			album = i.generateAlbumFromID3(ctx, albumAnalysis)
		} else {
			return err
		}
	} else {
		i.log.InfoContext(ctx, "album musicbrainz ID found", "mbID", albumAnalysis.AlbumID)
		album, err = i.generateAlbumFromMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			return err
		}
	}

	var artist types.Artist
	var tracks []types.Track

	if album.MusicBrainzID != nil {
		artist, err = i.generateArtistFromAlbumMbID(ctx, *album.MusicBrainzID)
		if err != nil {
			return err
		}
		tracks, err = i.generateTracksFromAlbumMbID(ctx, *album.MusicBrainzID)
		if err != nil {
			return err
		}
	} else {
		i.log.WarnContext(ctx, "no artist musicbrainz ID found, using ID3")
		artist = i.generateArtistFromID3(ctx, albumAnalysis)
	}

	albumFolderPath := utils.GenerateAlbumFolderPath(artist.Name, album.Name)

	for _, track := range albumAnalysis.Tracks {
		if track.File == "" {
			return fmt.Errorf("track '%s' does not have a file", *track.MbID)
		}

		tfp, err := utils.GenerateTrackPath(*track.DiscNumber, *track.TrackNumber, *track.Name, tag.FLAC, albumFolderPath)
		if err != nil {
			return err
		}

		if err = i.copyFile(ctx, track.File, filepath.Join(i.musicDir, tfp)); err != nil {
			return err
		}

		f, err := os.OpenFile(track.File, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		// FIXME: Do not hardcode Flac
		af, err := files.ParseAudioFile(f, tag.FLAC)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()

		if track.MbID != nil {
			for tidx := range tracks {
				if *tracks[tidx].MusicBrainzID == *track.MbID {
					// FIXME: Do not hardcode Flac
					tracks[tidx] = i.updateTrackData(tracks[tidx], af, tfp, tag.FLAC)
					break
				}
			}
		} else {
			t := types.Track{
				Title:       *track.Name,
				TrackNumber: track.TrackNumber,
				DiscNumber:  track.DiscNumber,
			}

			t = i.updateTrackData(t, af, tfp, tag.FLAC)
			tracks = append(tracks, t)
		}
	}

	tx, err := i.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingArtist *types.Artist

	artist.ID, existingArtist, err = i.addOrGetArtist(ctx, tx, artist)
	if err != nil {
		return err
	}

	if existingArtist != nil {
		if existingArtist.MusicBrainzID == nil && artist.MusicBrainzID != nil {
			err = tx.UpdateArtist(ctx, artist)
			if err != nil {
				return err
			}
		}
	}

	album.ArtistID = artist.ID

	album.ID, err = i.addOrGetAlbum(ctx, tx, album, artist.Name)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		track.AlbumID = &album.ID
		_, err := i.addOrUpdateTrack(ctx, tx, track)
		if err != nil {
			return err
		}
	}

	err = i.downloadAlbumCover(ctx, album, albumAnalysis.Covers)
	if err != nil {
		i.log.ErrorContext(ctx, "could not download album cover", "id", album.ID, "error", err.Error())
	}

	return tx.Commit()
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
	Files             []string
	Tracks            []Track
	Covers            []*tag.Picture
}

type Track struct {
	TrackNumber *int
	DiscNumber  *int
	Name        *string
	MbID        *string
	File        string
}

func (i Importer) generateAlbumFromMbID(ctx context.Context, mbID string) (types.Album, error) {
	mbAlbum, err := i.mb.GetAlbum(ctx, mbID)
	if err != nil {
		return types.Album{}, err
	}

	album := types.Album{
		MusicBrainzID: &mbAlbum.ID,
		Name:          mbAlbum.Title,
		SortName:      mbAlbum.Title,
		ReleaseDate:   &mbAlbum.Date.Time,
		Tracks:        &mbAlbum.TrackCount,
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
		Name:        albumName,
		SortName:    albumName,
		ReleaseDate: releaseDate,
		Tracks:      &totalTracks,
		Discs:       &totalDiscs,
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

	return artist, nil
}

func (i Importer) generateArtistFromID3(ctx context.Context, analysis AlbumAnalysisResults) types.Artist {
	artist := types.Artist{
		Name:     analysis.Id3Artist,
		SortName: analysis.Id3Artist,
	}

	return artist
}

func (i Importer) generateTracksFromAlbumMbID(ctx context.Context, albumMbID string) (tracks []types.Track, err error) {
	mbAlbum, err := i.mb.GetAlbum(ctx, albumMbID)
	if err != nil {
		return nil, err
	}

	for discIdx, media := range mbAlbum.Media {
		for _, mbTrack := range media.Tracks {
			tr := i.generateTrack(mbTrack, discIdx+1)
			tracks = append(tracks, tr)
		}
	}

	return tracks, nil
}

func (i Importer) generateTrack(mbTrack musicbrainz.Track, discNumber int) types.Track {
	return types.Track{
		Title:         mbTrack.Title,
		MusicBrainzID: &mbTrack.ID,
		TrackNumber:   &mbTrack.Position,
		DiscNumber:    &discNumber,
	}
}

func (i Importer) getID3Info(ctx context.Context, files []string) (artist, album string, releaseYear int, tracks []Track, pics []*tag.Picture, err error) {
	albums := []string{}
	artists := []string{}
	releaseYears := []int{}

	for _, f := range files {
		r, err := os.Open(f)
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
		track.File = f
		tn, _ := md.Track()
		dn, _ := md.Disc()

		if tn != 0 {
			track.TrackNumber = &tn
		}

		if dn != 0 {
			track.DiscNumber = &dn
		} else {
			track.DiscNumber = utils.Ptr(1)
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
