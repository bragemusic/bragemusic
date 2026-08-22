package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bragemusic/bragemusic/pkg/internalusers"
	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
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
		routes.New("PUT", "/me", s.updateProfile(), nil, routes.RouteMeta{
			Summary:             "Update profile data.",
			Description:         "Update profile data for the authenticated user.",
			ExpectedDescription: "Data updated",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("POST", "/me/img", s.uploadProfilePicture(), nil, routes.RouteMeta{
			Summary:             "Upload profile picture.",
			Description:         "Upload profile picture for the authenticated user.",
			ExpectedDescription: "Image uploaded",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("GET", "/me/tokens", s.usersListTokens(), nil, routes.RouteMeta{
			Summary:             "List tokens.",
			Description:         "List all tokens belonging to the logged in user.",
			ExpectedDescription: "Token data",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("POST", "/me/tokens/api", s.usersCreateApiToken(), nil, routes.RouteMeta{
			Summary:             "Create new API token.",
			Description:         "Create a new API token that will live forever.",
			ExpectedDescription: "Token data",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("DELETE", "/me/tokens/{tokenID}", s.usersDeleteToken(), nil, routes.RouteMeta{
			Summary:             "Delete token",
			Description:         "Delete a token with a specific ID.",
			ExpectedDescription: "Token deleted",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("POST", "/{userID}", s.updateUser(), []types.UserRole{types.UserRoleAdmin, types.UserRoleUsersUpdate}, routes.RouteMeta{
			Summary:             "Edit an existing user.",
			Description:         "Edit an exisiting user's information. Password is optional, every other pieces of information must be provided.",
			ExpectedDescription: "User updated",
			Tags:                []string{"Users"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("DELETE", "/{userID}", s.deleteUser(), []types.UserRole{types.UserRoleAdmin, types.UserRoleUsersDelete}, routes.RouteMeta{
			Summary:             "Delete an existing user.",
			Description:         "Delete an exisiting user from the server. This action cannot be reversed.",
			ExpectedDescription: "User deleted",
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

func (s *Server) deleteUser() routes.RouteFunc[ReqUsersBase, types.NoResponse] {
	return func(ctx context.Context, req ReqUsersBase, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.authPkg.RemoveUser(ctx, req.UserID); err != nil {
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

func (s *Server) updateUser() routes.RouteFunc[ReqUsersUpdate, types.NoResponse] {
	return func(ctx context.Context, req ReqUsersUpdate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.authPkg.UpdateUser(ctx, req.UserID, req.Email, req.Username, req.Password, req.Roles)
		if err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) updateProfile() routes.RouteFunc[ReqUsersUpdateProfile, types.NoResponse] {
	return func(ctx context.Context, req ReqUsersUpdateProfile, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.authPkg.UpdateProfile(ctx, user.ID, req.Email, req.Username, req.Password, req.NewPassword, req.NewPasswordConfirm)
		if err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) uploadProfilePicture() routes.RouteFunc[ReqNoContent, types.NoResponse] {
	return func(ctx context.Context, req ReqNoContent, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
		if err != nil {
			return resp, err
		}

		// Get the file from the form input "file"
		file, _, err := r.FormFile("file")
		if err != nil {
			return resp, err
		}
		defer file.Close()

		tempFolder, err := os.MkdirTemp(os.TempDir(), "brage-img")
		if err != nil {
			return resp, err
		}
		defer os.RemoveAll(tempFolder)

		orgImgPath := filepath.Join(tempFolder, fmt.Sprintf("%s.jpg", user.ID.String()))

		// Create the file on the server
		dst, err := os.Create(orgImgPath)
		if err != nil {
			return resp, err
		}
		defer dst.Close()

		// Copy the uploaded file's content to the destination file
		if _, err = io.Copy(dst, file); err != nil {
			return resp, err
		}

		if err = s.mediamgr.AddUserImage(ctx, orgImgPath, user.ID); err != nil {
			return resp, err
		}

		err = s.sseHub.SendToUser(user.ID, types.SSEServerMessage(types.SSEventMsgInfo, "Profile image Updated", "Your profile image has been updated"))
		if err != nil {
			s.log.WarnContext(ctx, "could not send server message", "error", err.Error())
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) usersCreateApiToken() routes.RouteFunc[ReqUsersApiTokenCreate, types.CreateUserApiTokenResp] {
	return func(ctx context.Context, req ReqUsersApiTokenCreate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.CreateUserApiTokenResp], err error) {
		token, _, err := s.authPkg.CreateAPIToken(ctx, req.Name, user.ID)
		if err != nil {
			return types.Response[types.CreateUserApiTokenResp]{}, err
		}

		err = s.sseHub.SendToUser(user.ID, types.SSEServerMessage(types.SSEventMsgInfo, "API-Token created", "A token was added to your account"))
		if err != nil {
			s.log.WarnContext(ctx, "could not send server message", "error", err.Error())
		}

		return types.Response[types.CreateUserApiTokenResp]{
			Payload: types.CreateUserApiTokenResp{
				Token: token,
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) usersDeleteToken() routes.RouteFunc[ReqUsersTokenDelete, types.NoResponse] {
	return func(ctx context.Context, req ReqUsersTokenDelete, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.authPkg.RemoveToken(ctx, req.TokenID, user.ID); err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		err = s.sseHub.SendToUser(user.ID, types.SSEServerMessage(types.SSEventMsgInfo, "Token deleted", "A token was successfully deleted"))
		if err != nil {
			s.log.WarnContext(ctx, "could not send server message", "error", err.Error())
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) usersListTokens() routes.RouteFunc[ReqNoContent, types.ListPayload[types.TokenLimited]] {
	return func(ctx context.Context, req ReqNoContent, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.TokenLimited]], err error) {
		tokens, err := s.authPkg.ListUserTokens(ctx, user.ID)
		if err != nil {
			return types.Response[types.ListPayload[types.TokenLimited]]{}, err
		}

		return types.Response[types.ListPayload[types.TokenLimited]]{
			Payload: types.ListPayload[types.TokenLimited]{
				Items: tokens,
				Count: len(tokens),
			},
			Status: http.StatusOK,
		}, nil
	}
}
