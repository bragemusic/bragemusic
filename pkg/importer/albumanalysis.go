package importer

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/samber/lo"
)

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

	matchedAlbum, err := i.getBestMatchedMbID(aids, id3Album)
	if err != nil {
		if errors.Is(err, ErrAlbumMbIDNotFound) {
			// FIXME return only ID3 info
			return AlbumAnalysisResults{}, errors.New("FIXME not implemented")
		}
		return AlbumAnalysisResults{}, err
	}

	mbAlbum, err := i.mb.GetAlbum(matchedAlbum.AlbumID)
	if err != nil {
		return AlbumAnalysisResults{}, err
	}

	aares := AlbumAnalysisResults{
		AlbumID:   mbAlbum.ID,
		Id3Artist: id3Artist,
		Id3Album:  id3Album,
	}

	for idx := range aids {
		id3Track := id3Tracks[idx]
		aid := aids[idx]
		file := files[idx]
		found := false

		for discNmbr, media := range mbAlbum.Media {
			for _, mbT := range media.Tracks {
				if lo.ContainsBy(aid, func(item acoustid.AcoustMatch) bool {
					return item.TrackID == mbT.ID
				}) {
					aares.Tracks = append(aares.Tracks, Track{
						TrackNumber: &mbT.Position,
						DiscNumber:  utils.Ptr(discNmbr + 1),
						Name:        &mbT.Title,
						MbID:        &mbT.ID,
						File:        file,
					})
					found = true
					break
				}
			}
		}

		if !found {
			aares.Tracks = append(aares.Tracks, Track{
				MbID:        nil,
				TrackNumber: id3Track.TrackNumber,
				DiscNumber:  id3Track.DiscNumber,
				Name:        id3Track.Name,
				File:        file,
			})
		}
	}

	for tidx := range aares.Tracks {
		if aares.Tracks[tidx].MbID == nil {
			if aares.Tracks[tidx].DiscNumber != nil && aares.Tracks[tidx].TrackNumber != nil {
				for discNmbr, media := range mbAlbum.Media {
					for _, mbT := range media.Tracks {
						if discNmbr == *aares.Tracks[tidx].DiscNumber && mbT.Position == *aares.Tracks[tidx].TrackNumber {
							fmt.Println("FoUND")
						}
					}
				}
			} else if aares.Tracks[tidx].Name != nil {
				for discNmbr, media := range mbAlbum.Media {
					for _, mbT := range media.Tracks {
						if utils.CompareTwoStrings(*aa, stringTwo string)
						fmt.Println(*aares.Tracks[tidx].Name)
					}
				}
			}

			// FIXME: match first on track and disc number. If not possible, match on name of track but make sure that the results are not already used by another track that is alredy matched.
			// If none of above works, only return ID3 tags as is
			fmt.Println("ejeje", *aares.Tracks[tidx].Name)
		}
	}

	fmt.Println(aares)
	_ = id3Artist
	_ = id3Tracks
	panic("hej")

	return AlbumAnalysisResults{}, nil
}

func (i Importer) getBestMatchedMbID(aids [][]acoustid.AcoustMatch, id3Album string) (MbAlbum, error) {
	mbAlbums := []MbAlbum{}

	for aIdx := range aids {
		for _, aid := range aids[aIdx] {
			mbAlbum := MbAlbum{
				AlbumID: aid.AlbumID,
			}

			for _, aIdObjs := range aids {
				match := slices.ContainsFunc(aIdObjs, func(ao acoustid.AcoustMatch) bool {
					return ao.AlbumID == aid.AlbumID
				})
				if match {
					mbAlbum.MatchedTracks++
					break
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
		return MbAlbum{}, ErrAlbumMbIDNotFound
	}

	return mbAlbums[0], nil
}
