package mediamanager

import (
	"context"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) SearchFull(ctx context.Context, searchTerm string) (si []types.SearchItem, err error) {
	if searchTerm == "" {
		return nil, nil
	}

	res, err := m.db.SearchFull(ctx, searchTerm, 10)
	if err != nil {
		return nil, m.berr.DatabaseError(err, types.EntitySearchItem, nil)
	}

	return res, nil
}

func (m MediaManager) UpdateSearchItems(ctx context.Context) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = tx.DeleteAllSearchItems(ctx); err != nil {
		return err
	}

	artists, err := tx.ListArtists(ctx, database.SortByName, database.SortAsc)
	if err != nil {
		return err
	}

	for _, a := range artists {
		si := types.SearchItem{
			Name:     a.Name,
			ID:       a.ID,
			Type:     types.EntityArtist,
			LinkID:   a.ID,
			LinkType: types.EntityArtist,
		}

		if err = tx.AddSearchItem(ctx, si); err != nil {
			return err
		}
	}

	albums, err := tx.ListAlbums(ctx)
	if err != nil {
		return err
	}

	for _, a := range albums {
		si := types.SearchItem{
			Name:     a.Name,
			ID:       a.ID,
			Type:     types.EntityAlbum,
			LinkID:   a.ID,
			LinkType: types.EntityAlbum,
		}

		if err = tx.AddSearchItem(ctx, si); err != nil {
			return err
		}
	}

	tracks, err := tx.ListTracks(ctx)
	if err != nil {
		return err
	}

	for _, t := range tracks {
		ats, err := tx.ListAlbumTracksByTrackID(ctx, t.ID)
		if err != nil {
			return err
		}

		for _, at := range ats {
			si := types.SearchItem{
				Name:     t.Title,
				ID:       t.ID,
				Type:     types.EntityTrack,
				LinkID:   at.AlbumID,
				LinkType: types.EntityAlbum,
			}

			if err := tx.AddSearchItem(ctx, si); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
