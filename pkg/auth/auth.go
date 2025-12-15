package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var (
	adminUUID        = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
	ErrTokenNotValid = errors.New("token not valid")
)

type Auth struct {
	hc  *HashCrypt
	db  database.AuthFace
	log *slog.Logger

	frontendTokenLongDuration  time.Duration
	frontendTokenShortDuration time.Duration
}

func (a Auth) CreateUser(ctx context.Context, userID uuid.UUID, email, username, password string, roles []types.UserRole) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user := types.User{
		ID:       userID,
		Email:    email,
		Username: username,
	}

	if _, err = tx.CreateUser(ctx, user); err != nil {
		return err
	}

	authIdentity := types.AuthIdentity{
		UserID:         userID,
		Provider:       types.AuthLocal,
		ProviderUserID: userID,
		Email:          email,
	}

	if _, err = tx.CreateAuthIdentity(ctx, authIdentity); err != nil {
		return err
	}

	pw, err := a.hc.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	localCredentials := types.LocalCredentials{
		UserID:       userID,
		PasswordHash: pw,
	}

	if err = tx.CreateLocalCredentials(ctx, localCredentials); err != nil {
		return err
	}

	for _, r := range roles {
		if err = tx.CreateUserScope(ctx, types.UserScope{UserID: userID, Role: r}); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (a Auth) UpdateUser(ctx context.Context, userID uuid.UUID, email, username string, password *string, roles []types.UserRole) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user := types.User{
		ID:       userID,
		Email:    email,
		Username: username,
	}

	if err = tx.UpdateUser(ctx, user); err != nil {
		return err
	}

	authIdentity, err := tx.GetAuthIdentityForUser(ctx, userID)
	if err != nil {
		return err
	}

	authIdentity.Email = email

	if err = tx.UpdateAuthIdentity(ctx, authIdentity); err != nil {
		return err
	}

	if password != nil {
		pw, err := a.hc.GenerateFromPassword(*password)
		if err != nil {
			return err
		}

		localCredentials, err := tx.GetLocalCredentialsForUser(ctx, userID)
		if err != nil {
			return err
		}

		localCredentials.PasswordHash = pw
		if err = tx.UpdateLocalCredentials(ctx, localCredentials); err != nil {
			return err
		}
	}

	currentRoles, err := tx.ListUserRoles(ctx, userID)
	if err != nil {
		return err
	}

	for _, r := range currentRoles {
		if !slices.Contains(roles, r) {
			if err = tx.RemoveUserScope(ctx, userID, r); err != nil {
				return err
			}
		}
	}

	for _, r := range roles {
		roleExists, err := tx.UserScopeExists(ctx, userID, r)
		if err != nil {
			return err
		}

		if !roleExists {
			if err = tx.CreateUserScope(ctx, types.UserScope{UserID: userID, Role: r}); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (a Auth) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = tx.RemoveUser(ctx, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (a Auth) SetAdmin(ctx context.Context, email, username, password string) error {
	userExists, err := a.db.UserExistsByID(ctx, adminUUID)
	if err != nil {
		return err
	}

	scope := []types.UserRole{types.UserRoleAdmin, types.UserRoleRead, types.UserRoleWrite}

	if userExists {
		if err = a.UpdateUser(ctx, adminUUID, email, username, &password, scope); err != nil {
			return err
		}
	} else {
		if err = a.CreateUser(ctx, adminUUID, email, username, password, scope); err != nil {
			return err
		}
	}

	return nil
}

func (a Auth) GenerateToken(ctx context.Context, userID uuid.UUID, tokenType types.TokenType, name *string) (token string, err error) {
	t := types.Token{
		UserID:    userID,
		Type:      tokenType,
		Name:      name,
		Hash:      "",
		ExpiresAt: nil,
	}

	switch tokenType {
	case types.TokenFrontendLong:
		expiresAt := time.Now().Add(a.frontendTokenLongDuration)
		t.ExpiresAt = &expiresAt
	case types.TokenFrontendShort:
		expiresAt := time.Now().Add(a.frontendTokenShortDuration)
		t.ExpiresAt = &expiresAt
	default:
		t.ExpiresAt = nil
	}

	token, err = generateToken()
	if err != nil {
		return "", err
	}

	t.Hash = hashToken(token)

	if _, err = a.db.CreateToken(ctx, t); err != nil {
		return "", err
	}

	return token, nil
}

func (a Auth) getUserFromTokenString(ctx context.Context, tokenString string) (types.UserDetails, error) {
	hash := hashToken(tokenString)

	token, err := a.db.GetTokenFromHash(ctx, hash)
	if err != nil {
		return types.UserDetails{}, ErrTokenNotValid
	}

	user, err := a.db.GetUserDetails(ctx, token.UserID)
	if err != nil {
		return types.UserDetails{}, err
	}

	return user, nil
}

func (a Auth) tokenFromHeader(ctx context.Context, authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("no Authorization header found")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Authorization header has wrong format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if !strings.HasPrefix(token, "brg_v1_") {
		return "", errors.New("token has wrong format")
	}

	return token, nil
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := a.tokenFromHeader(ctx, r.Header.Get("Authorization"))
		if err != nil {
			a.log.ErrorContext(ctx, err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user, err := a.getUserFromTokenString(ctx, token)
		if err != nil {
			a.log.ErrorContext(ctx, err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx = UpgradeContextWithUser(ctx, user)

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func New(db database.AuthFace, slogHandler slog.Handler) Auth {
	return Auth{
		hc:                         NewHashCrypt(),
		log:                        slog.New(slogHandler).With("service", "auth"),
		db:                         db,
		frontendTokenLongDuration:  7 * 24 * time.Hour,
		frontendTokenShortDuration: 24 * time.Hour,
	}
}
