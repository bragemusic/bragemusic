package authclient

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/adrg/xdg"
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

func (ac *AuthClient) LoadToken(ctx context.Context) error {
	path, err := ac.tokenPath()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	ac.sc.SetAuthToken(strings.TrimSpace(string(b)))

	user, err := ac.sc.GetUser(ctx)
	if err != nil {
		return err
	}

	for _, f := range ac.userCallbacks {
		f(&user)
	}

	return nil
}

func (ac *AuthClient) removeToken(ctx context.Context) error {
	ac.sc.SetAuthToken("")

	path, err := ac.tokenPath()
	if err != nil {
		return err
	}

	os.Remove(path)

	return nil
}

func (ac *AuthClient) saveToken(ctx context.Context, token string) error {
	ac.sc.SetAuthToken(token)

	path, err := ac.tokenPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(token), 0o600); err != nil {
		return err
	}

	return nil
}

func (ac *AuthClient) tokenPath() (string, error) {
	return xdg.StateFile("brage/token")
}

func New(sc *serverclient.ServerClient, slogHandler slog.Handler) AuthClient {
	return AuthClient{
		sc:  sc,
		log: slog.New(slogHandler).With("service", "auth-client"),
	}
}
