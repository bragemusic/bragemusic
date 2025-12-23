package authclient

import (
	"context"
	"log/slog"

	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
)

type AuthClient struct {
	sc            *serverclient.ServerClient
	log           *slog.Logger
	userCallbacks []func(*types.UserDetails)
}

func (s *AuthClient) RegisterUserCallback(f func(*types.UserDetails)) {
	s.userCallbacks = append(s.userCallbacks, f)
}

func (ac *AuthClient) Login(ctx context.Context, username, password string, longLivedToken bool) error {
	loginResp, err := ac.sc.Login(ctx, username, password, longLivedToken)
	if err != nil {
		return err
	}

	if err = ac.saveToken(ctx, loginResp.Token); err != nil {
		return err
	}

	user, err := ac.sc.GetUser(ctx)
	if err != nil {
		return err
	}

	for _, f := range ac.userCallbacks {
		f(&user)
	}

	return nil
}

func (ac *AuthClient) Logout(ctx context.Context) error {
	if err := ac.removeToken(ctx); err != nil {
		return err
	}

	for _, f := range ac.userCallbacks {
		f(nil)
	}

	return nil
}

func (ac *AuthClient) removeToken(ctx context.Context) error {
	ac.sc.SetAuthToken("")
	ac.log.WarnContext(ctx, "now token should be removed. TODO")
	return nil
}

func (ac *AuthClient) saveToken(ctx context.Context, token string) error {
	ac.sc.SetAuthToken(token)
	ac.log.WarnContext(ctx, "now token should be saved. TODO")
	return nil
}

func New(sc *serverclient.ServerClient, slogHandler slog.Handler) AuthClient {
	return AuthClient{
		sc:  sc,
		log: slog.New(slogHandler).With("service", "auth-client"),
	}
}
