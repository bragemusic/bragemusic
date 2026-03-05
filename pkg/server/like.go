package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) likeRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{likeID}", s.getLike(), nil, routes.RouteMeta{
			Summary:             "Get a like by ID",
			Description:         "Returns metadata for the specified like.",
			ExpectedDescription: "Like metadata",
			Tags:                []string{"Likes"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getLike() routes.RouteFunc[ReqLikesGet, types.Like] {
	return func(ctx context.Context, req ReqLikesGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.Like], err error) {
		like, err := s.mediamgr.GetLike(ctx, req.LikeID, user.ID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.Like]{
			Payload: like,
			Status:  http.StatusOK,
		}, nil
	}
}
