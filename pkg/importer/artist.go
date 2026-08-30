package importer

import (
	"context"
	"errors"

	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (i Importer) generateArtist(ctx context.Context, albumAnalysis AlbumAnalysisResults, userID uuid.UUID) (artist types.Artist, err error) {
	if albumAnalysis.AlbumID != "" {
		artist, err = i.generateArtistFromAlbumMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			// Issue #140: If mbID is not found for the artist, use ID3 instead
			if errors.Is(err, ErrArtistMbIDNotFound) {
				i.log.WarnContext(ctx, "no artist musicbrainz ID found, using ID3")
				artist = i.generateArtistFromID3(ctx, albumAnalysis)
				return artist, nil
			}
			return types.Artist{}, err
		}
	} else {
		i.log.WarnContext(ctx, "no artist musicbrainz ID found, using ID3")
		artist = i.generateArtistFromID3(ctx, albumAnalysis)
	}
	return artist, nil
}

func (i Importer) addArtist(ctx context.Context, tx database.DatabaseFace, artist types.Artist, albumAnalysis AlbumAnalysisResults, userID uuid.UUID) (artistID uuid.UUID, err error) {
	artistID, existingArtist, err := i.addOrGetArtist(ctx, tx, artist, userID)
	if err != nil {
		return uuid.Nil, err
	}

	// Issue #140: Should fix artists with empty string as mbID
	if artist.MusicBrainzID != nil && *artist.MusicBrainzID == "" {
		artist.MusicBrainzID = nil
	}

	if existingArtist != nil {
		artist.ID = existingArtist.ID
		if existingArtist.MusicBrainzID == nil && artist.MusicBrainzID != nil {
			err = tx.UpdateArtist(ctx, artist, userID)
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	return artistID, nil
}
