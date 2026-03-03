package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) ratingRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{ratingID}", s.getRating(), nil, routes.RouteMeta{
			Summary:             "Get a rating by ID",
			Description:         "Returns metadata for the specified rating.",
			ExpectedDescription: "Rating metadata",
			Tags:                []string{"Ratings"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getRating() routes.RouteFunc[ReqRatingsGet, types.Rating] {
	return func(ctx context.Context, req ReqRatingsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.Rating], err error) {
		rating, err := s.mediamgr.GetRating(ctx, req.RatingID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.Rating]{
			Payload: rating,
			Status:  http.StatusOK,
		}, nil
	}
}
