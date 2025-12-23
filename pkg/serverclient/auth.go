package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/server"
)

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
