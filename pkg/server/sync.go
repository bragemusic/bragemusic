package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) syncRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("POST", "/", s.sync(), nil, routes.RouteMeta{
			Summary:             "Get sync state since time",
			Description:         "Returns all events since the specified time, tied to the user. Used for syncing client and server.",
			ExpectedDescription: "Sync data",
			Tags:                []string{"Sync"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("POST", "/play-history", s.syncPlayHistory(), nil, routes.RouteMeta{
			Summary:             "Sync play history items between client and server",
			Description:         "Adds play history items given in the request and returns all other new items",
			ExpectedDescription: "Play history data",
			Tags:                []string{"Sync"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) sync() routes.RouteFunc[ReqSyncSync, types.SyncState] {
	return func(ctx context.Context, req ReqSyncSync, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.SyncState], err error) {
		syncState, err := s.mediamgr.GetSyncState(ctx, req.ChangesSince)
		if err != nil {
			return resp, err
		}

		return types.Response[types.SyncState]{
			Payload: syncState,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) syncPlayHistory() routes.RouteFunc[ReqSyncPlayHistory, types.PlayHistorySyncState] {
	return func(ctx context.Context, req ReqSyncPlayHistory, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.PlayHistorySyncState], err error) {
		syncState, err := s.mediamgr.SyncPlayHistory(ctx, req.ChangesSince, req.UpdatedClientItems)
		if err != nil {
			return resp, err
		}

		return types.Response[types.PlayHistorySyncState]{
			Payload: syncState,
			Status:  http.StatusOK,
		}, nil
	}
}
