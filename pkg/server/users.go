package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/internalusers"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) userRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listUsers(), []types.UserRole{types.UserRoleAdmin, types.UserRoleUsersRead}, routes.RouteMeta{
			Summary:             "List all users.",
			Description:         "List all users on the server, including backend machine users.",
			ExpectedDescription: "List of users",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),

		routes.New("POST", "/", s.createUser(), []types.UserRole{types.UserRoleAdmin, types.UserRoleUsersCreate}, routes.RouteMeta{
			Summary:             "Create a new user.",
			Description:         "Create a new user with a local auth provider",
			ExpectedDescription: "User created",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
	}
}

func (s *Server) createUser() routes.RouteFunc[ReqUsersCreate, types.NoResponse] {
	return func(ctx context.Context, req ReqUsersCreate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.authPkg.CreateUser(ctx, req.Email, req.Username, req.Password, req.Roles)
		if err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) listUsers() routes.RouteFunc[ReqUsersList, types.ListPayload[types.UserDetails]] {
	return func(ctx context.Context, req ReqUsersList, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.UserDetails]], err error) {
		users, err := s.authPkg.ListUsers(ctx)
		if err != nil {
			return resp, err
		}

		if req.IncludeMachineUsers {
			users = append(users, internalusers.GetIntenalUsers()...)
		}

		if req.Count {
			return types.Response[types.ListPayload[types.UserDetails]]{
				Payload: types.ListPayload[types.UserDetails]{
					Items: nil,
					Count: len(users),
				},
				Status: http.StatusOK,
			}, nil
		}

		return types.Response[types.ListPayload[types.UserDetails]]{
			Payload: types.ListPayload[types.UserDetails]{
				Items: users,
				Count: len(users),
			},
			Status: http.StatusOK,
		}, nil
	}
}
