package importer

import (
	"cmp"
	"context"
	"errors"
	"slices"

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/musicbrainz"
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
	id3Artist, id3Album, id3Year, id3Tracks, mdPics, err := i.getID3Info(ctx, files)
	if err != nil {
		if errors.Is(err, ErrId3AlbumNotFound) {
			return AlbumAnalysisResults{}, errors.New("FIXME: See if there is something we can do when no ID3Album exists.. Maybe we can continue with one of the MB albums anyway and it should be ok. Exile on main street can be ued for test")
		}
		return AlbumAnalysisResults{}, err
	}

	matchedAlbum, err := i.getBestMatchedMbID(aids, id3Album)
	if err != nil {
		if errors.Is(err, ErrAlbumMbIDNotFound) {
			i.log.WarnContext(ctx, "could not find MusicBrainzID, using ID3 instead", "album", id3Album)
			return AlbumAnalysisResults{
				Id3Artist:      id3Artist,
				Id3Album:       id3Album,
				Id3ReleaseDate: id3Year,
				Files:          files,
				Tracks:         id3Tracks,
				Covers:         mdPics,
			}, ErrAlbumMbIDNotFound
		}
		return AlbumAnalysisResults{}, err
	}

	mbAlbum, err := i.mb.GetAlbum(ctx, matchedAlbum.AlbumID)
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

	aares.Tracks, err = i.matchRemainingTracks(ctx, aares.Tracks, mbAlbum)
	if err != nil {
		return AlbumAnalysisResults{}, err
	}

	aares.Covers = mdPics

	return aares, nil
}

func (i Importer) matchRemainingTracks(ctx context.Context, tracks []Track, mbAlbum musicbrainz.Release) ([]Track, error) {
	// List used MusicBrainz IDs
	usedIDs := []string{}
	for _, t := range tracks {
		if t.MbID != nil {
			usedIDs = append(usedIDs, *t.MbID)
		}
	}

	availableTracks := []Track{}

	// Create tracks for all remaining MusicBrainz IDs
	for discNmbr, media := range mbAlbum.Media {
		for _, mbT := range media.Tracks {
			if slices.ContainsFunc(usedIDs, func(tID string) bool {
				return mbT.ID == tID
			}) {
				continue
			}

			availableTracks = append(availableTracks, Track{
				TrackNumber: &mbT.Position,
				DiscNumber:  utils.Ptr(discNmbr + 1),
				Name:        &mbT.Title,
				MbID:        &mbT.ID,
			})
		}
	}

	// Try to match the remaining tracks
	for tidx := range tracks {
		if tracks[tidx].MbID == nil {
			var err error
			tracks[tidx], availableTracks, err = i.matchTrack(ctx, tracks[tidx], availableTracks)
			if err != nil {
				return nil, err
			}
		}
	}

	return tracks, nil
}

func (i Importer) matchTrack(ctx context.Context, track Track, availableTracks []Track) (Track, []Track, error) {
	type match struct {
		V   float32
		Idx int
	}

	// First try matching using disc and track number
	if track.DiscNumber != nil && track.TrackNumber != nil {
		for tidx, t := range availableTracks {
			if *t.DiscNumber == *track.DiscNumber && *t.TrackNumber == *track.TrackNumber {
				i.log.DebugContext(ctx, "matched track using disc and track number", "track_id", *t.MbID)
				t.File = track.File
				availableTracks = slices.Delete(availableTracks, tidx, tidx)
				return t, availableTracks, nil
			}
		}
	}

	stringMatch := []match{}

	// If not possible with track and disc number, compare names of the tracks
	for tidx, t := range availableTracks {
		stringMatch = append(stringMatch, match{V: utils.CompareTwoStrings(*track.Name, *t.Name), Idx: tidx})
	}

	// Sort based on best match
	slices.SortFunc(stringMatch, func(a, b match,
	) int {
		return cmp.Compare(b.V, a.V)
	})

	if len(stringMatch) == 0 {
		return Track{}, availableTracks, errors.New("no tracks available")
	}

	tidx := stringMatch[0].Idx
	t := availableTracks[tidx]

	i.log.DebugContext(ctx, "matched track using name", "track_id", *t.MbID, "match_score", stringMatch[0].V)

	t.File = track.File
	availableTracks = slices.Delete(availableTracks, tidx, tidx)

	return t, availableTracks, nil
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
			// create own compare function to find earliest official album
		)
	})

	if len(mbAlbums) == 0 {
		return MbAlbum{}, ErrAlbumMbIDNotFound
	}

	return mbAlbums[0], nil
}
