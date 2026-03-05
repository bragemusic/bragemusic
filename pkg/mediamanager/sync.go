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

	artists, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityArtist, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, artists...)

	albums, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityAlbum, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, albums...)

	tracks, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityTrack, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, tracks...)

	albumArtists, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityAlbumArtist, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, albumArtists...)

	albumTracks, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityAlbumTrack, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, albumTracks...)

	playlists, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityPlaylist, since, &user.ID)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, playlists...)

	playlistTracks, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityPlaylistTrack, since, &user.ID)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, playlistTracks...)

	mediaFiles, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityMediaFile, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, mediaFiles...)

	ratings, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityRating, since, nil)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, ratings...)

	likes, err := m.db.ListEntityEventsByEntityType(ctx, types.EntityLike, since, &user.ID)
	if err != nil {
		return types.SyncState{}, err
	}
	st.New = append(st.New, likes...)

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
