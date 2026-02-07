package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
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

	ac.sc.SetAuthToken(loginResp.Token)

	user, err := ac.sc.GetUser(ctx)
	if err != nil {
		return err
	}

	if err = ac.saveToken(ctx, loginResp.Token, user.ID); err != nil {
		return err
	}

	if err = ac.saveUserDetails(ctx, user); err != nil {
		return err
	}

	for _, f := range ac.userCallbacks {
		f(&user)
	}

	return nil
}

func (ac *AuthClient) LogoutServerUser(ctx context.Context, userID uuid.UUID) error {
	if err := ac.removeToken(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (ac *AuthClient) LoginLocalUser(ctx context.Context, userID uuid.UUID) error {
	path, err := ac.tokenPath(userID)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	ac.sc.SetAuthToken(strings.TrimSpace(string(b)))

	userPath, err := ac.userDetailsPath(userID)
	if err != nil {
		return err
	}

	var user types.UserDetails

	bUser, err := os.ReadFile(userPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			user, err = ac.sc.GetUser(ctx)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if err = json.Unmarshal(bUser, &user); err != nil {
			return err
		}
	}

	for _, f := range ac.userCallbacks {
		f(&user)
	}

	return nil
}

func (ac *AuthClient) GetCachedUsers(ctx context.Context) (users []types.UserDetails, err error) {
	path, err := ac.usersPath()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			userID, err := uuid.FromString(entry.Name())
			if err != nil {
				continue
			}

			userPath, err := ac.userDetailsPath(userID)
			if err != nil {
				continue
			}

			var user types.UserDetails

			bUser, err := os.ReadFile(userPath)
			if err != nil {
				continue
			}

			if err = json.Unmarshal(bUser, &user); err != nil {
				continue
			}

			users = append(users, user)
		}
	}

	return users, nil
}

func (ac *AuthClient) removeToken(ctx context.Context, userID uuid.UUID) error {
	ac.sc.SetAuthToken("")

	path, err := ac.tokenPath(userID)
	if err != nil {
		return err
	}

	os.Remove(path)

	return nil
}

func (ac *AuthClient) saveToken(ctx context.Context, token string, userID uuid.UUID) error {
	path, err := ac.tokenPath(userID)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(token), 0o600); err != nil {
		fmt.Println("kaka")
		return err
	}

	return nil
}

func (ac *AuthClient) saveUserDetails(ctx context.Context, user types.UserDetails) error {
	path, err := ac.userDetailsPath(user.ID)
	if err != nil {
		return err
	}

	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, b, 0o600); err != nil {
		return err
	}

	return nil
}

func (ac *AuthClient) usersPath() (string, error) {
	return xdg.StateFile(filepath.Join("brage", "users"))
}

func (ac *AuthClient) tokenPath(userID uuid.UUID) (string, error) {
	return xdg.StateFile(filepath.Join("brage", "users", userID.String(), "token"))
}

func (ac *AuthClient) userDetailsPath(userID uuid.UUID) (string, error) {
	return xdg.StateFile(filepath.Join("brage", "users", userID.String(), "user.json"))
}

func New(sc *serverclient.ServerClient, slogHandler slog.Handler) AuthClient {
	return AuthClient{
		sc:  sc,
		log: slog.New(slogHandler).With("service", "auth-client"),
	}
}
