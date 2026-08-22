package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) Login(ctx context.Context, email, password string, longLivedToken bool) (resp types.LoginResp, err error) {
	u, err := url.JoinPath(s.baseUrl, "auth", "login")
	if err != nil {
		return types.LoginResp{}, err
	}

	body := types.LoginReq{
		Email:          email,
		Password:       password,
		LongLivedToken: longLivedToken,
	}

	if err := s.doPostJson(ctx, u, body, &resp); err != nil {
		return types.LoginResp{}, err
	}

	return resp, nil
}

func (s ServerClient) Logout(ctx context.Context) (err error) {
	u, err := url.JoinPath(s.baseUrl, "auth", "logout")
	if err != nil {
		return err
	}

	if err := s.doGetJson(ctx, u, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DeleteToken(ctx context.Context, tokenID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "auth", "tokens", tokenID.String())
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) ListUserRoles(ctx context.Context) ([]types.UserRole, error) {
	u, err := url.JoinPath(s.baseUrl, "auth", "user-roles")
	if err != nil {
		return nil, err
	}

	resp := types.ListPayload[types.UserRole]{}

	if err := s.doGetJson(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}
