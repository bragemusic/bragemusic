package importer

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bragemusic/core/pkg/acoustid"
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

	files := []string{}
	for _, f := range osFiles {
		files = append(files, filepath.Join(folder, f.Name()))
	}

	var album types.Album

	// returnera ett objekt h'rist'llet med filnamn pa varje track tillsammans med mbid o lite id3 (namn, nummer osv), och albumid. Sa blir det lattare att mathc asen
	albumMbID, trackMbIDs, err := i.getAlbumMbID(ctx, files)
	if err != nil {
		if errors.Is(err, ErrAlbumMbIDNotFound) {
			return errors.New("not implemented to add non MB album. TODO: Do something else")
		} else {
			return err
		}
	} else {
		i.log.InfoContext(ctx, "album musicbrainz ID found", "mbID", albumMbID)
		album, err = i.generateAlbumFromMbID(ctx, albumMbID)
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
		return errors.New("non mb ID artist not implemented")
	}

	albumFolderPath := utils.GenerateAlbumFolderPath(artist.Name, album.Name)

	fmt.Println(trackMbIDs)
	for idx, filename := range files {
		trackMbId := trackMbIDs[idx]
		if trackMbId != nil {
			for tidx, t := range tracks {
				fmt.Println(*tracks[tidx].MusicBrainzID, *trackMbId)
				if *tracks[tidx].MusicBrainzID == *trackMbId {
					tfp, err := utils.GenerateTrackPath(*t.DiscNumber, *t.TrackNumber, t.Title, tag.FLAC, albumFolderPath)
					if err != nil {
						return err
					}
					tracks[tidx].FilePath = tfp
					// har ska filen kopieras till ratt stalle. Man borde kanske ocksa tabort anvanda filer och trackMbIds ur listorna
					_ = filename
					break
				}
			}
		} else {
			i.log.WarnContext(ctx, "non mb ID track not implemented")
		}
	}

	fmt.Println(album)
	fmt.Println(artist)
	fmt.Println(tracks)

	return nil
}

type MbAlbum struct {
	AlbumID           string
	MbTrackIDs        []*string
	MatchedTracks     int
	ID3AlbumNameMatch float32
}

type AlbumAnalysisResults struct {
	AlbumID           string
	Id3Artist         string
	Id3Album          string
	MatchedTracks     int
	ID3AlbumNameMatch float32
	Files             []string
	FilesMbIds        []*string
	Tracks            []Track
}

type Track struct {
	TrackNumber *int
	DiscNumber  *int
	Name        *string
	MbID        *string
}

func (i Importer) generateAlbumFromMbID(ctx context.Context, mbID string) (types.Album, error) {
	mbAlbum, err := i.mb.GetAlbum(mbID)
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

func (i Importer) generateArtistFromAlbumMbID(ctx context.Context, mbID string) (types.Artist, error) {
	mbAlbum, err := i.mb.GetAlbum(mbID)
	if err != nil {
		return types.Artist{}, err
	}

	if len(mbAlbum.ArtistCredit) == 0 {
		return types.Artist{}, ErrArtistMbIDNotFound
	}

	artistMbID := mbAlbum.ArtistCredit[0].Artist.ID

	mbArtist, err := i.mb.GetArtist(artistMbID)
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

func (i Importer) generateTracksFromAlbumMbID(ctx context.Context, albumMbID string) (tracks []types.Track, err error) {
	mbAlbum, err := i.mb.GetAlbum(albumMbID)
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

func (i Importer) analyzeAlbum(ctx context.Context, files []string) (AlbumAnalysisResults, error) {
	aids := [][]acoustid.AcoustMatch{}

	for _, f := range files {
		aid, err := i.aid.GetMusicBrainzAlbumID(f)
		if err != nil {
			return AlbumAnalysisResults{}, err
		}

		aids = append(aids, aid)
	}

	i.log.InfoContext(ctx, "getting album name from ID3")
	id3Artist, id3Album, id3Tracks, err := i.getID3Info(ctx, files)
	if err != nil {
		return AlbumAnalysisResults{}, err
	}

	mbAlbums := []MbAlbum{}

	for aIdx := range aids {
		for _, aid := range aids[aIdx] {
			mbAlbum := AlbumAnalysisResults{
				AlbumID: aid.AlbumID,
			}

			for _, aIdObjs := range aids {
				match := slices.ContainsFunc(aIdObjs, func(ao acoustid.AcoustMatch) bool {
					return ao.AlbumID == aid.AlbumID
				})
				if match {
					mbAlbum.MatchedTracks++
					for _, ao := range aIdObjs {
						if ao.AlbumID == aid.AlbumID {
							mbAlbum.Tracks = append(mbAlbum.Tracks, Track{
								TrackNumber: tracks[iiii],
								DiscNumber:  new(int),
								Name:        new(string),
								MbID:        &ao.TrackID,
							})
							break
						}
					}
				} else {
					mbAlbum.MbTrackIDs = append(mbAlbum.MbTrackIDs, nil)
				}
			}

			mbAlbum.ID3AlbumNameMatch = utils.CompareTwoStrings(aid.AlbumName, id3Album)

			if !slices.ContainsFunc(mbAlbums, func(ma MbAlbum) bool {
				return ma.AlbumID == mbAlbum.AlbumID
			}) {
				mbAlbums = append(mbAlbums, mbAlbum)
			}
		}
	}

	slices.SortFunc(mbAlbums, func(a, b MbAlbum) int {
		return cmp.Or(
			cmp.Compare(b.MatchedTracks, a.MatchedTracks),
			cmp.Compare(b.ID3AlbumNameMatch, a.ID3AlbumNameMatch),
		)
	})

	if len(mbAlbums) == 0 {
		return AlbumAnalysisResults{}, ErrAlbumMbIDNotFound
	}

	return AlbumAnalysisResults{}, ErrAlbumMbIDNotFound
	// return mbAlbums[0].AlbumID, mbAlbums[0].MbTrackIDs, nil
}

// func (i Importer) getAlbumMbID(ctx context.Context, files []string) (string, []*string, error) {
// 	aids := [][]acoustid.AcoustMatch{}

// 	for _, f := range files {
// 		aid, err := i.aid.GetMusicBrainzAlbumID(f)
// 		if err != nil {
// 			return "", nil, err
// 		}

// 		aids = append(aids, aid)
// 	}

// 	i.log.InfoContext(ctx, "getting album name from ID3")
// 	id3Album, err := i.getAlbumNameID3(ctx, files)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	mbAlbums := []MbAlbum{}

// 	for aIdx := range aids {
// 		for _, aid := range aids[aIdx] {
// 			mbAlbum := MbAlbum{
// 				AlbumID: aid.AlbumID,
// 			}

// 			for _, aIdObjs := range aids {
// 				match := slices.ContainsFunc(aIdObjs, func(ao acoustid.AcoustMatch) bool {
// 					return ao.AlbumID == aid.AlbumID
// 				})
// 				if match {
// 					mbAlbum.MatchedTracks++
// 					for _, ao := range aIdObjs {
// 						if ao.AlbumID == aid.AlbumID {
// 							mbAlbum.MbTrackIDs = append(mbAlbum.MbTrackIDs, &ao.TrackID)
// 							break
// 						}
// 					}
// 				} else {
// 					mbAlbum.MbTrackIDs = append(mbAlbum.MbTrackIDs, nil)
// 				}
// 			}

// 			mbAlbum.ID3AlbumNameMatch = utils.CompareTwoStrings(aid.AlbumName, id3Album)

// 			if !slices.ContainsFunc(mbAlbums, func(ma MbAlbum) bool {
// 				return ma.AlbumID == mbAlbum.AlbumID
// 			}) {
// 				mbAlbums = append(mbAlbums, mbAlbum)
// 			}
// 		}
// 	}

// 	slices.SortFunc(mbAlbums, func(a, b MbAlbum) int {
// 		return cmp.Or(
// 			cmp.Compare(b.MatchedTracks, a.MatchedTracks),
// 			cmp.Compare(b.ID3AlbumNameMatch, a.ID3AlbumNameMatch),
// 		)
// 	})

// 	if len(mbAlbums) == 0 {
// 		return "", nil, ErrAlbumMbIDNotFound
// 	}

// 	return mbAlbums[0].AlbumID, mbAlbums[0].MbTrackIDs, nil
// }

func (i Importer) getID3Info(ctx context.Context, files []string) (artist, album string, tracks []Track, err error) {
	albums := []string{}
	artists := []string{}

	for _, f := range files {
		r, err := os.Open(f)
		if err != nil {
			return "", "", nil, err
		}

		md, err := tag.ReadFrom(r)
		if err != nil {
			return "", "", nil, err
		}

		if md.Album() != "" {
			albums = append(albums, md.Album())
		}

		if md.Artist() != "" {
			artists = append(artists, md.Artist())
		}

		track := Id3Track{}
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

		tracks = append(tracks, track)

		r.Close()
	}

	i.log.DebugContext(ctx, "found album names in ID3", "names", albums)

	if len(albums) == 0 {
		return "", "", nil, ErrId3AlbumNotFound
	}

	return utils.HighestCount(artists), utils.HighestCount(albums), tracks, nil
}
