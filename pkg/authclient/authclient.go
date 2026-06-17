package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var ErrNoCachedToken = errors.New("no cached token")

type AuthClient struct {
	sc            *serverclient.ServerClient
	log           *slog.Logger
	userCallbacks []func(*types.UserDetails)
}

func (ac *AuthClient) ServerClient() *serverclient.ServerClient {
	return ac.sc
}

func (ac *AuthClient) RegisterUserCallback(f func(*types.UserDetails)) {
	ac.userCallbacks = append(ac.userCallbacks, f)
}

func (ac *AuthClient) UserCallback(user *types.UserDetails) {
	for _, f := range ac.userCallbacks {
		f(user)
	}
}

func (ac *AuthClient) LoginCachedServerUser(ctx context.Context, tokenType types.TokenType) error {
	path, err := ac.tokenPath(tokenType)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoCachedToken
		}
		return err
	}

	ac.sc.SetAuthToken(strings.TrimSpace(string(b)))

	return nil
}

func (ac *AuthClient) Login(ctx context.Context, username, password string, tokenType types.TokenType) (types.UserDetails, error) {
	var token string

	loginResp, err := ac.sc.Login(ctx, username, password, tokenType == types.TokenFrontendLong)
	if err != nil {
		return types.UserDetails{}, err
	}

	ac.sc.SetAuthToken(loginResp.Token)

	if tokenType == types.TokenAPI {
		token, err = ac.sc.CreateAPIToken(ctx, "Desktop Client")
		if err != nil {
			return types.UserDetails{}, err
		}
	} else {
		token = loginResp.Token
	}

	user, err := ac.sc.GetUser(ctx)
	if err != nil {
		return types.UserDetails{}, err
	}

	if err = ac.saveToken(ctx, token, user.ID, tokenType); err != nil {
		return types.UserDetails{}, err
	}

	if err = ac.saveUserDetails(ctx, user); err != nil {
		return types.UserDetails{}, err
	}

	return user, err
}

func (ac *AuthClient) LoginToken(ctx context.Context, token string) (types.UserDetails, error) {
	ac.sc.SetAuthToken(token)

	user, err := ac.sc.GetUser(ctx)
	if err != nil {
		return types.UserDetails{}, err
	}

	return user, err
}

func (ac *AuthClient) LogoutServerUser(ctx context.Context, tokenType types.TokenType) error {
	if err := ac.removeToken(ctx, tokenType); err != nil {
		return err
	}

	return nil
}

// func (ac *AuthClient) LoginLocalUser(ctx context.Context, userID uuid.UUID, runCallback bool) (types.UserDetails, error) {
// 	path, err := ac.tokenPath()
// 	if err != nil {
// 		return types.UserDetails{}, err
// 	}

// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		if !errors.Is(err, os.ErrNotExist) {
// 			return types.UserDetails{}, err
// 		}
// 	}

// 	ac.sc.SetAuthToken(strings.TrimSpace(string(b)))

// 	userPath, err := ac.userDetailsPath(userID)
// 	if err != nil {
// 		return types.UserDetails{}, err
// 	}

// 	var user types.UserDetails

// 	bUser, err := os.ReadFile(userPath)
// 	if err != nil {
// 		if errors.Is(err, os.ErrNotExist) {
// 			user, err = ac.sc.GetUser(ctx)
// 			if err != nil {
// 				return types.UserDetails{}, err
// 			}
// 		} else {
// 			return types.UserDetails{}, err
// 		}
// 	} else {
// 		if err = json.Unmarshal(bUser, &user); err != nil {
// 			return types.UserDetails{}, err
// 		}
// 	}

// 	if runCallback {
// 		ac.UserCallback(&user)
// 	}

// 	return user, nil
// }

func (ac *AuthClient) GetCachedUser(ctx context.Context) (user *types.UserDetails, err error) {
	path, err := ac.userPath()
	if err != nil {
		return nil, err
	}

	bUser, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	if err = json.Unmarshal(bUser, &user); err != nil {
		return nil, err
	}

	return user, nil
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

func (ac *AuthClient) RemoveToken(ctx context.Context, id uuid.UUID) error {
	return ac.sc.DeleteToken(ctx, id)
}

func (ac *AuthClient) removeToken(ctx context.Context, tokenType types.TokenType) error {
	ac.sc.SetAuthToken("")

	path, err := ac.tokenPath(tokenType)
	if err != nil {
		return err
	}

	os.Remove(path)

	return nil
}

func (ac *AuthClient) saveToken(ctx context.Context, token string, userID uuid.UUID, tokenType types.TokenType) error {
	path, err := ac.tokenPath(tokenType)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(token), 0o600); err != nil {
		return err
	}

	return nil
}

func (ac *AuthClient) saveUserDetails(ctx context.Context, user types.UserDetails) error {
	path, err := ac.userPath()
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

func (ac *AuthClient) userPath() (string, error) {
	return xdg.StateFile(filepath.Join("brage", "user-sync.json"))
}

func (ac *AuthClient) tokenPath(tokenType types.TokenType) (string, error) {
	if tokenType == types.TokenAPI {
		return xdg.StateFile(filepath.Join("brage", "token-sync"))
	}
	return xdg.StateFile(filepath.Join("brage", "token-stream"))
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
