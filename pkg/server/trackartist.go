package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) trackArtistRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{trackArtistID}", s.getTrackArtistByID(), nil, routes.RouteMeta{
			Summary:             "Retrieve a track artist by ID.",
			Description:         "Returns metadata about the specified track artist.",
			ExpectedDescription: "Metadata about the track artist",
			Tags:                []string{"Track Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getTrackArtistByID() routes.RouteFunc[ReqTrackArtistsGet, types.TrackArtist] {
	return func(ctx context.Context, req ReqTrackArtistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.TrackArtist], err error) {
		trackArtist, err := s.mediamgr.GetTrackArtistByID(ctx, req.TrackArtistID)
		if err != nil {
			return types.Response[types.TrackArtist]{}, err
		}

		return types.Response[types.TrackArtist]{
			Payload: trackArtist,
			Status:  http.StatusOK,
		}, nil
	}
}
