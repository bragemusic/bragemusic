package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
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
