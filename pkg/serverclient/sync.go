package serverclient

import (
	"context"
	"net/url"
	"time"

	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) GetSyncState(ctx context.Context, changesSince time.Time) (syncState types.SyncState, err error) {
	u, err := url.JoinPath(s.baseUrl, "sync")
	if err != nil {
		return types.SyncState{}, err
	}

	payload := server.SyncReq{ChangesSince: changesSince}

	if err := s.doPostJson(ctx, u, payload, &syncState); err != nil {
		return types.SyncState{}, err
	}

	return syncState, nil
}

func (s ServerClient) SyncPlayHistory(ctx context.Context, changesSince time.Time, newItems []types.PlayHistory) (syncState types.PlayHistorySyncState, err error) {
	u, err := url.JoinPath(s.baseUrl, "/sync/play-history")
	if err != nil {
		return types.PlayHistorySyncState{}, err
	}

	payload := server.SyncPlayHistoryReq{
		ChangesSince:       changesSince,
		UpdatedClientItems: newItems,
	}

	if err := s.doPostJson(ctx, u, payload, &syncState); err != nil {
		return types.PlayHistorySyncState{}, err
	}

	return syncState, nil
}
