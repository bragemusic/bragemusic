package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) playlistTrackRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{playlistTrackID}", s.getPlaylistTrack(), nil, routes.RouteMeta{
			Summary:             "Get a playlist track by ID",
			Description:         "Returns metadata for the specified playlist track.",
			ExpectedDescription: "Playlist track metadata",
			Tags:                []string{"Playlist Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("DELETE", "/{playlistTrackID}", s.deletePlaylistTrack(), nil, routes.RouteMeta{
			Summary:             "Delete a playlist track by ID",
			Description:         "Deletes the specified playlist track.",
			ExpectedDescription: "Playlist track successfully deleted",
			Tags:                []string{"Playlist Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
	}
}

func (s *Server) getPlaylistTrack() routes.RouteFunc[ReqPlaylistTracksGet, types.PlaylistTrack] {
	return func(ctx context.Context, req ReqPlaylistTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.PlaylistTrack], err error) {
		pt, err := s.mediamgr.GetPlaylistTrack(ctx, req.PlaylistTrackID, user.ID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.PlaylistTrack]{
			Payload: pt,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) deletePlaylistTrack() routes.RouteFunc[ReqPlaylistTracksGet, types.NoResponse] {
	return func(ctx context.Context, req ReqPlaylistTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.mediamgr.DeletePlaylistTrack(ctx, req.PlaylistTrackID, user.ID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}
