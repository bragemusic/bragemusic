package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var adminUUID = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))

type Auth struct {
	db  database.DatabaseFace
	log *slog.Logger
}

func (a Auth) SetAdmin(ctx context.Context, email, username, password string) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userExists, err := tx.UserExistsByID(ctx, adminUUID)
	if err != nil {
		return err
	}

	user := types.User{
		ID:       adminUUID,
		Email:    email,
		Username: username,
	}

	if userExists {
		if err = tx.UpdateUser(ctx, user); err != nil {
			return err
		}
	} else {
		if _, err = tx.CreateUser(ctx, user); err != nil {
			return err
		}
	}

	authIdentity, err := tx.GetAuthIdentityForUser(ctx, adminUUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		authIdentity = types.AuthIdentity{
			UserID:         adminUUID,
			Provider:       types.AuthLocal,
			ProviderUserID: adminUUID,
			Email:          email,
		}

		if _, err = tx.CreateAuthIdentity(ctx, authIdentity); err != nil {
			return err
		}
	} else {
		authIdentity.Email = email
		if err = tx.UpdateAuthIdentity(ctx, authIdentity); err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func New(db database.DatabaseFace, slogHandler slog.Handler) Auth {
	return Auth{
		log: slog.New(slogHandler).With("service", "auth"),
		db:  db,
	}
}
