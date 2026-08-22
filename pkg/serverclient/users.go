package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) CreateUser(ctx context.Context, email, username, password string, roles []types.UserRole) error {
	u, err := url.JoinPath(s.baseUrl, "api", "users")
	if err != nil {
		return err
	}

	req := types.CreateUserReq{
		Email:    email,
		Username: username,
		Password: password,
		Roles:    roles,
	}

	if err := s.doPostJson(ctx, u, req, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) CreateAPIToken(ctx context.Context, name string) (token string, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "users", "me", "tokens", "api")
	if err != nil {
		return "", err
	}

	req := types.CreateUserApiTokenReq{
		Name: name,
	}

	resp := types.CreateUserApiTokenResp{}

	if err := s.doPostJson(ctx, u, req, &resp); err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (s ServerClient) DeleteUser(ctx context.Context, id uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "users", id.String())
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DeleteUserToken(ctx context.Context, tokenID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "users", "me", "tokens", tokenID.String())
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) ListUsers(ctx context.Context, includeMachineUsers bool) (users []types.UserDetails, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "users")
	if err != nil {
		return nil, err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	q := ur.Query()
	if includeMachineUsers {
		q.Set("machineUsers", "true")
	}

	ur.RawQuery = q.Encode()

	resp := types.ListPayload[types.UserDetails]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (s ServerClient) ListUserTokens(ctx context.Context) (tokens []types.TokenLimited, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "users", "me", "tokens")
	if err != nil {
		return nil, err
	}

	resp := types.ListPayload[types.TokenLimited]{}

	if err := s.doGetJson(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (s ServerClient) UpdateUser(ctx context.Context, userID uuid.UUID, email, username string, password *string, roles []types.UserRole) error {
	u, err := url.JoinPath(s.baseUrl, "api", "users", userID.String())
	if err != nil {
		return err
	}

	req := types.UpdateUserReq{
		Email:    email,
		Username: username,
		Password: password,
		Roles:    roles,
	}

	if err := s.doPostJson(ctx, u, req, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) UpdateProfile(ctx context.Context, data types.UpdateProfileReq) error {
	u, err := url.JoinPath(s.baseUrl, "api", "users", "me")
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, data, nil); err != nil {
		return err
	}

	return nil
}
