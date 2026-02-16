package mediamanager

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) GetSyncState(ctx context.Context, since time.Time) (st types.SyncState, err error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return types.SyncState{}, err
	}

	st.Time = time.Now()

	st.CreatedOrUpdated.Playlists, err = m.db.ListUpdatedPlaylists(ctx, since, user.ID)
	if err != nil {
		return types.SyncState{}, err
	}

	st.CreatedOrUpdated.PlaylistTracks, err = m.db.ListUpdatedPlaylistTracks(ctx, since, user.ID)
	if err != nil {
		return types.SyncState{}, err
	}

	st.CreatedOrUpdated.MediaFiles, err = m.db.ListUpdatedMediaFiles(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

	// st.Deleted.AlbumArtists, err = m.db.ListEntityEvents(ctx, types.EntityEventDelete, types.EntityAlbumArtist, since)
	// if err != nil {
	// 	return types.SyncState{}, err
	// }

	st.Deleted.Playlists, err = m.db.ListEntityEvents(ctx, types.EntityEventDelete, types.EntityPlaylist, since)
	if err != nil {
		return types.SyncState{}, err
	}

	st.Deleted.PlaylistTracks, err = m.db.ListEntityEvents(ctx, types.EntityEventDelete, types.EntityPlaylistTrack, since)
	if err != nil {
		return types.SyncState{}, err
	}

	artists, err := m.db.ListEntityEventsTemp(ctx, types.EntityArtist, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, artists...)

	albums, err := m.db.ListEntityEventsTemp(ctx, types.EntityAlbum, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, albums...)

	tracks, err := m.db.ListEntityEventsTemp(ctx, types.EntityTrack, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, tracks...)

	albumArtists, err := m.db.ListEntityEventsTemp(ctx, types.EntityAlbumArtist, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, albumArtists...)

	albumTracks, err := m.db.ListEntityEventsTemp(ctx, types.EntityAlbumTrack, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, albumTracks...)

	ratings, err := m.db.ListEntityEventsTemp(ctx, types.EntityRating, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, ratings...)

	return
}

func (m MediaManager) SyncPlayHistory(ctx context.Context, since time.Time, newItems []types.PlayHistory) (st types.PlayHistorySyncState, err error) {
	st.Time = time.Now()

	phs, err := m.db.ListUpdatedPlayHistory(ctx, since)
	if err != nil {
		return types.PlayHistorySyncState{}, err
	}

	st.RemoteItems = phs

	if len(newItems) == 0 {
		m.log.DebugContext(ctx, fmt.Sprintf("no new play history items given by client, %d found in server", len(phs)))
		return st, nil
	}

	tx, err := m.db.Begin(ctx)
	defer tx.Rollback()

	for _, item := range newItems {
		if _, err := m.db.AddPlayHistoryStruct(ctx, item); err != nil {
			return types.PlayHistorySyncState{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return types.PlayHistorySyncState{}, err
	}

	return
}
