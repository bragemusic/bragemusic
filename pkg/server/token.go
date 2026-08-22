package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) tokenRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("DELETE", "/{tokenID}", s.deleteToken(), nil, routes.RouteMeta{
			Summary:             "Delete a token by ID.",
			Description:         "Deletes a token owned by the logged in user. All sessions using the selected token will lose auth.",
			ExpectedDescription: "Token deleted",
			Tags:                []string{"Token"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
	}
}

func (s *Server) deleteToken() routes.RouteFunc[ReqTokensBase, types.NoResponse] {
	return func(ctx context.Context, req ReqTokensBase, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.authPkg.RemoveToken(ctx, req.TokenID, user.ID); err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}
