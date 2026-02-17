package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) ListUsers(ctx context.Context) (users []types.User, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "users")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s ServerClient) Login(ctx context.Context, email, password string, longLivedToken bool) (resp server.LoginResp, err error) {
	u, err := url.JoinPath(s.baseUrl, "auth", "login")
	if err != nil {
		return server.LoginResp{}, err
	}

	body := server.LoginReq{
		Email:          email,
		Password:       password,
		LongLivedToken: longLivedToken,
	}

	if err := s.doPostJson(ctx, u, body, &resp); err != nil {
		return server.LoginResp{}, err
	}

	return resp, nil
}
