package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) albumTrackRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{albumTrackID}", s.getAlbumTrackByID(), nil, routes.RouteMeta{
			Summary:             "Retrieve an album track by ID.",
			Description:         "Returns metadata about the specified album track.",
			ExpectedDescription: "Metadata about the album track",
			Tags:                []string{"Album Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getAlbumTrackByID() routes.RouteFunc[ReqAlbumTracksGet, types.AlbumTrack] {
	return func(ctx context.Context, req ReqAlbumTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.AlbumTrack], err error) {
		albumTrack, err := s.mediamgr.GetAlbumTrackByID(ctx, req.AlbumTrackID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.AlbumTrack]{
			Payload: albumTrack,
			Status:  http.StatusOK,
		}, nil
	}
}
