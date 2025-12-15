package auth

import (
	"context"
	"log/slog"
	"slices"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var adminUUID = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))

type Auth struct {
	hc  *HashCrypt
	db  database.DatabaseFace
	log *slog.Logger
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

func New(db database.DatabaseFace, slogHandler slog.Handler) Auth {
	return Auth{
		hc:  NewHashCrypt(),
		log: slog.New(slogHandler).With("service", "auth"),
		db:  db,
	}
}
