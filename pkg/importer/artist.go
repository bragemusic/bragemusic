package importer

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (i Importer) addArtist(ctx context.Context, tx database.DatabaseFace, albumAnalysis AlbumAnalysisResults, userID uuid.UUID) (artistID uuid.UUID, err error) {
	var artist types.Artist

	if albumAnalysis.AlbumID != "" {
		artist, err = i.generateArtistFromAlbumMbID(ctx, albumAnalysis.AlbumID)
		if err != nil {
			return uuid.Nil, err
		}
	} else {
		i.log.WarnContext(ctx, "no artist musicbrainz ID found, using ID3")
		artist = i.generateArtistFromID3(ctx, albumAnalysis)
	}

	artistID, existingArtist, err := i.addOrGetArtist(ctx, tx, artist, userID)
	if err != nil {
		return uuid.Nil, err
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
