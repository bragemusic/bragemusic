package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) userRoleRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listUserRoles(), []types.UserRole{types.UserRoleAdmin}, routes.RouteMeta{
			Summary:             "List all user roles.",
			Description:         "Lists all available user roles in the system.",
			ExpectedDescription: "List of user roles",
			Tags:                []string{"User Roles"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) listUserRoles() routes.RouteFunc[ReqNoContent, types.ListPayload[types.UserRole]] {
	return func(ctx context.Context, req ReqNoContent, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.UserRole]], err error) {
		return types.Response[types.ListPayload[types.UserRole]]{
			Payload: types.ListPayload[types.UserRole]{
				Items: types.AllUserRoles,
				Count: len(types.AllUserRoles),
			},
			Status: http.StatusOK,
		}, nil
	}
}
