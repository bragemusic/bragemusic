package mediamanager

import (
	"context"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/types"
)

func (m MediaManager) GetSyncState(ctx context.Context, since time.Time) (st types.SyncState, err error) {
	st.Time = time.Now()

	st.Artists, err = m.db.ListUpdatedArtists(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

	st.Albums, err = m.db.ListUpdatedAlbums(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

	st.Tracks, err = m.db.ListUpdatedTracks(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

	st.AlbumArtists, err = m.db.ListUpdatedAlbumArtists(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

	st.AlbumTracks, err = m.db.ListUpdatedAlbumTracks(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

	st.MediaFiles, err = m.db.ListUpdatedMediaFiles(ctx, since)
	if err != nil {
		return types.SyncState{}, err
	}

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
