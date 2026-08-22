package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/bragemusic/bragemusic/pkg/bragerr"
	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var (
	adminUUID             = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
	ErrTokenNotValid      = errors.New("token not valid")
	ErrInvalidCredentials = errors.New("invalid user credentials")
	ErrNoAuthHeader       = errors.New("no authorization header")
	ErrNoAuthCookie       = errors.New("no auth cookie")
	ErrInvalidScheme      = errors.New("invalid authorization scheme")
	ErrInvalidToken       = errors.New("invalid token format")
)

type Auth struct {
	hc   *HashCrypt
	db   database.AuthFace
	log  *slog.Logger
	berr bragerr.BragErrFactory

	frontendTokenLongDuration  time.Duration
	frontendTokenShortDuration time.Duration
}

func (a Auth) GetUserFromContext(ctx context.Context) (types.UserDetails, error) {
	user, err := UserFromContext(ctx)
	return user, err
}

func (a Auth) ListUsers(ctx context.Context) ([]types.UserDetails, error) {
	users, err := a.db.ListUsers(ctx)
	return users, err
}

func (a Auth) CreateUser(ctx context.Context, email, username, password string, roles []types.UserRole) error {
	userID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	return a.createUser(ctx, userID, email, username, password, roles)
}

func (a Auth) createUser(ctx context.Context, userID uuid.UUID, email, username, password string, roles []types.UserRole) error {
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

func (a Auth) UpdateProfile(ctx context.Context, userID uuid.UUID, email, username, password, newPassword, newPasswordConfirm *string) error {
	user, err := a.db.GetUserDetails(ctx, userID)
	if err != nil {
		return err
	}

	var uEmail, uUsername string

	if email == nil {
		uEmail = user.Email
	} else {
		uEmail = *email
	}

	if username == nil {
		uUsername = user.Username
	} else {
		uUsername = *username
	}

	if newPassword != nil {
		if newPassword == nil || newPasswordConfirm == nil {
			return a.berr.Unauthenticated(errors.New("new passwords are nil"))
		}

		if password == nil {
			return a.berr.Unauthenticated(errors.New("password is nil"))
		}

		if *newPassword != *newPasswordConfirm {
			return a.berr.Unauthenticated(errors.New("new passwords are not same"))
		}

		localCreds, err := a.db.GetLocalCredentialsForUser(ctx, user.ID)
		if err != nil {
			return ErrInvalidCredentials
		}

		passMatch, err := a.hc.ComparePasswordAndHash(*password, localCreds.PasswordHash)
		if err != nil {
			return ErrInvalidCredentials
		}

		if !passMatch {
			return ErrInvalidCredentials
		}
	}

	return a.UpdateUser(ctx, userID, uEmail, uUsername, newPassword, user.Roles)
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
		if err = a.createUser(ctx, adminUUID, email, username, password, scope); err != nil {
			return err
		}
	}

	return nil
}

func (a Auth) ListUserTokens(ctx context.Context, userID uuid.UUID) ([]types.TokenLimited, error) {
	tokens, err := a.db.ListUserTokens(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (a Auth) CreateLoginToken(ctx context.Context, email, password string, longLivedToken bool) (token string, expiresIn int, err error) {
	user, err := a.db.GetUserFromEmail(ctx, email)
	if err != nil {
		return "", 0, a.berr.ErrInvalidUserCreds(err)
	}

	localCreds, err := a.db.GetLocalCredentialsForUser(ctx, user.ID)
	if err != nil {
		return "", 0, a.berr.ErrInvalidUserCreds(err)
	}

	passMatch, err := a.hc.ComparePasswordAndHash(password, localCreds.PasswordHash)
	if err != nil {
		return "", 0, a.berr.ErrInvalidUserCreds(err)
	}

	if !passMatch {
		return "", 0, a.berr.ErrInvalidUserCreds(errors.New("password does not match"))
	}

	tokenType := types.TokenFrontendShort
	if longLivedToken {
		tokenType = types.TokenFrontendLong
	}

	token, expiresIn, err = a.generateToken(ctx, user.ID, tokenType, nil)
	if err != nil {
		return "", 0, a.berr.ErrInvalidUserCreds(err)
	}

	return token, expiresIn, nil
}

func (a Auth) CreateAPIToken(ctx context.Context, name string, userID uuid.UUID) (token string, expiresIn int, err error) {
	token, expiresIn, err = a.generateToken(ctx, userID, types.TokenAPI, &name)
	if err != nil {
		return "", 0, ErrInvalidCredentials
	}

	return token, expiresIn, nil
}

func (a Auth) RemoveToken(ctx context.Context, tokenID, userID uuid.UUID) error {
	token, err := a.db.GetToken(ctx, tokenID)
	if err != nil {
		return a.berr.DatabaseError(err, types.EntityToken, &tokenID)
	}

	// NOTE: In the future we can add a better check. Maybe admins should be able to delete any token in the system
	if token.UserID != userID {
		return a.berr.ItemAccessDenied(errors.New("user is not permitted to delete token"), types.EntityToken, tokenID)
	}

	if err := a.db.RemoveToken(ctx, tokenID); err != nil {
		return a.berr.DatabaseError(err, types.EntityToken, &tokenID)
	}

	return nil
}

func (a Auth) TokenCleanupJob(ctx context.Context) error {
	a.log.InfoContext(ctx, "started token cleanup job")

	tokensDeleted, err := a.db.RemoveExpiredTokens(ctx)
	if err != nil {
		return err
	}

	a.log.InfoContext(ctx, "token cleanup job finsished", "tokens_deleted", tokensDeleted)

	return nil
}

func (a Auth) generateToken(ctx context.Context, userID uuid.UUID, tokenType types.TokenType, name *string) (token string, expiresIn int, err error) {
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
	case types.TokenAPI:
		t.ExpiresAt = nil
	default:
		return "", 0, errors.New("forbidden token type")
	}

	token, err = generateToken()
	if err != nil {
		return "", 0, err
	}

	t.Hash = hashToken(token)

	if _, err = a.db.CreateToken(ctx, t); err != nil {
		return "", 0, err
	}

	if t.ExpiresAt != nil {
		expiresIn = int(t.ExpiresAt.Sub(time.Now()).Seconds())
	}

	return token, expiresIn, nil
}

func (a Auth) getUserFromTokenString(ctx context.Context, tokenString string) (user types.UserDetails, tokenID uuid.UUID, err error) {
	hash := hashToken(tokenString)

	token, err := a.db.GetTokenFromHash(ctx, hash)
	if err != nil {
		return types.UserDetails{}, uuid.Nil, ErrTokenNotValid
	}

	if err = a.validateToken(ctx, token); err != nil {
		return types.UserDetails{}, uuid.Nil, ErrTokenNotValid
	}

	if err = a.db.UpdateTokenLastUsed(ctx, token.ID); err != nil {
		return types.UserDetails{}, uuid.Nil, err
	}

	user, err = a.db.GetUserDetails(ctx, token.UserID)
	if err != nil {
		return types.UserDetails{}, uuid.Nil, err
	}

	return user, token.ID, nil
}

func (a Auth) validateToken(ctx context.Context, token types.Token) error {
	switch token.Type {
	case types.TokenAPI:
		return nil
	default:
		if token.ExpiresAt == nil {
			return ErrTokenNotValid
		}
	}

	if token.ExpiresAt.Before(time.Now()) {
		return ErrTokenNotValid
	}

	return nil
}

func (a Auth) tokenFromCookie(ctx context.Context, cookieName string, r *http.Request) (string, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", ErrNoAuthCookie
		}
		return "", err
	}

	return c.Value, nil
}

func (a Auth) tokenFromHeader(ctx context.Context, authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrInvalidScheme
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if !strings.HasPrefix(token, "brg_v1_") {
		return "", ErrInvalidToken
	}

	return token, nil
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, headerErr := a.tokenFromHeader(ctx, r.Header.Get("Authorization"))
		if headerErr != nil {
			if errors.Is(headerErr, ErrNoAuthHeader) {
				var cookieErr error
				token, cookieErr = a.tokenFromCookie(ctx, "brage_session_token", r)
				if cookieErr != nil {
					bragerr.HandleHttpResponse(ctx, a.berr.Unauthenticated(cookieErr), w, a.log)
					return
				}
			} else {
				bragerr.HandleHttpResponse(ctx, a.berr.Unauthenticated(headerErr), w, a.log)
				return
			}
		}

		user, tokenID, err := a.getUserFromTokenString(ctx, token)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, a.berr.Unauthenticated(err), w, a.log)
			return
		}

		ctx = UpgradeContextWithUser(ctx, user)
		ctx = UpgradeContextWithTokenID(ctx, tokenID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a Auth) RoleCheckMiddleware(roles ...types.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, err := a.GetUserFromContext(ctx)
			if err != nil {
				bragerr.HandleHttpResponse(ctx, a.berr.Unauthenticated(err), w, a.log)
			}

			for _, role := range roles {
				if user.HasRole(role) {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			bragerr.HandleHttpResponse(
				ctx,
				a.berr.Unauthenticated(errors.New("user does not have one of the required roles")).With("roles", roles),
				w,
				a.log,
			)
		})
	}
}

func New(db database.AuthFace, slogHandler slog.Handler) Auth {
	return Auth{
		hc:                         NewHashCrypt(),
		log:                        slog.New(slogHandler).With("service", "auth"),
		db:                         db,
		frontendTokenLongDuration:  7 * 24 * time.Hour,
		frontendTokenShortDuration: 24 * time.Hour,
		berr:                       bragerr.NewFactory("auth"),
	}
}
