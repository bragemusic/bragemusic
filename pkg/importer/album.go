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

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/files"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
	"github.com/gofrs/uuid/v5"
)

var (
	ErrAlbumMbIDNotFound  = errors.New("could not find MusicBrainz ID for album")
	ErrArtistMbIDNotFound = errors.New("could not find MusicBrainz ID for artist")
	ErrId3AlbumNotFound   = errors.New("could not get album name from ID3")
	ErrAlbumNotFound      = errors.New("no existing album found")
)

func (i Importer) addAlbum(ctx context.Context, tx database.DatabaseFace, albumAnalysis AlbumAnalysisResults, existingAlbum types.Album) (album types.Album, err error) {
	if albumAnalysis.AlbumID != "" {
		album, err = i.generateAlbumFromMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			return types.Album{}, err
		}
	} else {
		album = i.generateAlbumFromID3(ctx, albumAnalysis)
	}

	if existingAlbum.ID != uuid.Nil {
		// Exisiting album is not created from musicbrainz id but new one is, update it
		if existingAlbum.MusicBrainzID == nil && album.MusicBrainzID != nil {
			album.ID = existingAlbum.ID
			if err = tx.UpdateAlbum(ctx, album); err != nil {
				return types.Album{}, err
			}
		} else {
			album = existingAlbum
		}
	} else {
		album.ID, err = tx.AddAlbum(ctx, album)
		if err != nil {
			return types.Album{}, err
		}
	}

	return album, nil
}

func (i Importer) importAlbumFilesOld(ctx context.Context, folder string) error {
	osFiles, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	filenames := []string{}
	for _, f := range osFiles {
		filenames = append(filenames, filepath.Join(folder, f.Name()))
	}

	var album types.Album

	albumAnalysis, err := i.analyzeAlbum(ctx, []types.MediaFile{})
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
	var albumTracks []types.AlbumTrack

	if album.MusicBrainzID != nil {
		artist, err = i.generateArtistFromAlbumMbID(ctx, *album.MusicBrainzID)
		if err != nil {
			return err
		}
		tracks, albumTracks, err = i.generateTracksFromAlbumMbID(ctx, *album.MusicBrainzID)
		if err != nil {
			return err
		}
	} else {
		i.log.WarnContext(ctx, "no artist musicbrainz ID found, using ID3")
		artist = i.generateArtistFromID3(ctx, albumAnalysis)
	}

	tx, err := i.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, track := range albumAnalysis.Tracks {
		// if track.File == "" {
		if "" == "" {
			return fmt.Errorf("track '%s' does not have a file", *track.MbID)
		}

		// f, err := os.OpenFile(track.File, os.O_RDONLY, os.ModePerm)
		f, err := os.OpenFile("", os.O_RDONLY, os.ModePerm)
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

		// checksum, err := utils.FileSHA256(track.File)
		checksum, err := utils.FileSHA256("")
		if err != nil {
			return err
		}

		mediafile := types.MediaFile{
			DurationMs: af.DurationMS(),
			Bitrate:    af.Bitrate(),
			SampleRate: af.SampleRate(),
			FileSize:   af.FileSize(),
			// FIXME: Do not hardcode Flac
			Codec:    types.CodecFlac,
			Checksum: checksum,
		}

		mfId, err := tx.AddMediaFile(ctx, mediafile)
		if err != nil {
			return err
		}

		mffp := filepath.Join(i.musicDir, fmt.Sprintf("%s.%s", mfId.String(), mediafile.Codec))

		if err = i.copyFile(ctx, "", mffp); err != nil {
			// if err = i.copyFile(ctx, track.File, mffp); err != nil {
			return err
		}

		if track.MbID != nil {
			for tidx := range tracks {
				if *tracks[tidx].MusicBrainzID == *track.MbID {
					tracks[tidx].MediaFile = utils.Ptr(mfId)
					break
				}
			}
		} else {
			t := types.Track{
				Title:     *track.Name,
				MediaFile: utils.Ptr(mfId),
			}
			tracks = append(tracks, t)
			albumTracks = append(albumTracks, types.AlbumTrack{
				DiscNumber:  *track.DiscNumber,
				TrackNumber: *track.TrackNumber,
			})
		}
	}

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

	// FIXME: new tables
	// album.ArtistID = artist.ID

	album.ID, err = i.addOrGetAlbum(ctx, tx, album, artist.Name)
	if err != nil {
		return err
	}

	// FIXME: create album artist link

	fmt.Println(len(tracks), len(albumTracks))
	for idx := range tracks {
		// track.AlbumID = &album.ID
		trackID, newTrack, err := i.addOrUpdateTrack(ctx, tx, tracks[idx], album.ID)
		if err != nil {
			return err
		}
		albumTracks[idx].AlbumID = album.ID
		albumTracks[idx].TrackID = trackID

		if newTrack {
			if err = tx.AddAlbumTrack(ctx, albumTracks[idx]); err != nil {
				return err
			}
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

		if tn != 0 {
			track.TrackNumber = &tn
		} else {
			track.TrackNumber = utils.Ptr(1)
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

func (i Importer) getExistingAlbum(ctx context.Context, tx database.DatabaseFace, albumAnalysis AlbumAnalysisResults) (types.Album, error) {
	if albumAnalysis.AlbumID != "" {
		album, err := tx.GetAlbumFromMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return types.Album{}, err
			}
		} else {
			return album, nil
		}
	}

	album, err := tx.GetAlbumFromArtistAndName(ctx, albumAnalysis.Id3Artist, albumAnalysis.Id3Album)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return types.Album{}, err
		}
	} else {
		return album, nil
	}

	return types.Album{}, ErrAlbumNotFound
}
