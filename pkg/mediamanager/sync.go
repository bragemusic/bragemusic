package mediamanager

import (
	"context"
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

	return
}
