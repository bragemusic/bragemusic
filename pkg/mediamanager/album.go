package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

func (m MediaManager) GetAlbum(ctx context.Context, albumID uuid.UUID) (types.Album, error) {
	album, err := m.db.GetAlbumFromID(ctx, albumID)
	if err != nil {
		return types.Album{}, m.berr.DatabaseError(err, types.EntityAlbum, &albumID)
	}

	return album, nil
}

func (m MediaManager) GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (types.AlbumDetailed, error) {
	album, err := m.db.GetAlbumDetailed(ctx, albumID)
	if err != nil {
		return types.AlbumDetailed{}, m.berr.DatabaseError(err, types.EntityAlbum, &albumID)
	}

	return album, nil
}

func (m MediaManager) ListAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error) {
	albums, err := m.db.ListAlbumsByArtist(ctx, artistID, sortBy, sortOrder)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityAlbum, &artistID)
	}

	return albums, nil
}

func (m MediaManager) ListFeaturedAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error) {
	albums, err := m.db.ListFeaturedAlbumsByArtist(ctx, artistID, sortBy, sortOrder)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityAlbum, &artistID)
	}

	return albums, nil
}

func (m MediaManager) ListAlbums(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) (albums []types.AlbumDetailed, err error) {
	artists, err := m.db.ListArtists(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntityArtist, nil)
	}

	for _, artist := range artists {
		alb, err := m.ListAlbumsByArtist(ctx, artist.ID, sortBy, sortOrder)
		if err != nil {
			return nil, err
		}

		for _, album := range alb {
			_, idx, exists := lo.FindIndexOf(albums, func(item types.AlbumDetailed) bool {
				return item.ID == album.ID
			})

			if !exists {
				album.ArtistNames = append(album.ArtistNames, artist.Name)
				album.ArtistIDs = append(album.ArtistIDs, artist.ID.String())
				albums = append(albums, album)
				continue
			}

			albums[idx].ArtistNames = append(albums[idx].ArtistNames, artist.Name)
			albums[idx].ArtistIDs = append(albums[idx].ArtistIDs, artist.ID.String())

		}

	}

	return albums, nil
}

func (m MediaManager) GetAlbumArtist(ctx context.Context, albumID, artistID uuid.UUID, role types.ArtistRole) (types.AlbumArtist, error) {
	albumArtist, err := m.db.GetAlbumArtist(ctx, albumID, artistID, types.ArtistRole(role))
	if err != nil {
		return types.AlbumArtist{}, m.berr.DatabaseError(err, types.EntityAlbumArtist, &albumID)
	}

	return albumArtist, nil
}

func (m MediaManager) GetAlbumArtistByID(ctx context.Context, id uuid.UUID) (types.AlbumArtist, error) {
	albumArtist, err := m.db.GetAlbumArtistByID(ctx, id)
	if err != nil {
		return types.AlbumArtist{}, m.berr.DatabaseError(err, types.EntityAlbumArtist, &id)
	}

	return albumArtist, nil
}

func (m MediaManager) GetAlbumTrack(ctx context.Context, albumID uuid.UUID, discNumber, trackNumber int) (types.AlbumTrack, error) {
	albumArtist, err := m.db.GetAlbumTrack(ctx, albumID, discNumber, trackNumber)
	if err != nil {
		return types.AlbumTrack{}, err
	}

	return albumArtist, nil
}

func (m MediaManager) GetAlbumTrackByID(ctx context.Context, id uuid.UUID) (types.AlbumTrack, error) {
	albumArtist, err := m.db.GetAlbumTrackByID(ctx, id)
	if err != nil {
		return types.AlbumTrack{}, m.berr.DatabaseError(err, types.EntityAlbumTrack, &id)
	}

	return albumArtist, nil
}

func (m MediaManager) UpdateAlbum(ctx context.Context, albumID uuid.UUID, albumData types.AlbumUpdate, userID uuid.UUID) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existingAlbum, err := tx.GetAlbumFromID(ctx, albumID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityAlbum, &albumID)
	}

	albumData.Owner = existingAlbum.Owner

	album := types.Album{
		ID:        albumID,
		AlbumBase: albumData.AlbumBase,
		CreatedAt: existingAlbum.CreatedAt,
		UpdatedAt: existingAlbum.UpdatedAt,
	}

	err = tx.UpdateAlbum(ctx, album, userID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityAlbum, &albumID)
	}

	existingAlbumArtists, err := tx.ListAlbumArtistsByAlbumID(ctx, albumID)
	if err != nil {
		return m.berr.DatabaseError(err, types.EntityAlbumArtist, nil)
	}

	for _, artistID := range albumData.Artists {
		exists := lo.ContainsBy(existingAlbumArtists, func(item types.AlbumArtist) bool {
			return item.ArtistID == artistID
		})

		if !exists {
			aa := types.AlbumArtist{
				AlbumID:  albumID,
				ArtistID: artistID,
				Role:     types.ArPrimary,
				// FIXME: Need to do something about positions
				Position: 0,
			}
			if _, err := tx.AddAlbumArtist(ctx, aa, userID); err != nil {
				return m.berr.DatabaseError(err, types.EntityAlbumArtist, nil)
			}
		}
	}

	for _, existingAa := range existingAlbumArtists {
		exists := lo.ContainsBy(albumData.Artists, func(item uuid.UUID) bool {
			return item == existingAa.ArtistID
		})

		if !exists {
			if err := tx.DeleteAlbumArtist(ctx, existingAa.ID, userID); err != nil {
				return m.berr.DatabaseError(err, types.EntityAlbumArtist, &existingAa.ID)
			}
		}
	}

	return tx.Commit()
}

func (m MediaManager) CountAlbums(ctx context.Context) (int, error) {
	cnt, err := m.db.CountAlbums(ctx)
	if err != nil {
		return 0, m.berr.DatabaseError(err, types.EntityAlbum, nil)
	}

	return cnt, nil
}
