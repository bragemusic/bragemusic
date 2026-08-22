package auth

import (
	"context"
	"errors"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type (
	userContextKeyType    struct{}
	tokenIDContextKeyType struct{}
)

var (
	ucKey = userContextKeyType{}
	tcKey = tokenIDContextKeyType{}
)

func TokenIDFromContext(ctx context.Context) (uuid.UUID, error) {
	td, ok := ctx.Value(tcKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("no token in context")
	}

	return td, nil
}

func UserFromContext(ctx context.Context) (types.UserDetails, error) {
	ud, ok := ctx.Value(ucKey).(types.UserDetails)
	if !ok {
		return types.UserDetails{}, errors.New("no user in context")
	}

	return ud, nil
}

func UpgradeContextWithTokenID(ctx context.Context, tokenID uuid.UUID) context.Context {
	return context.WithValue(ctx, tcKey, tokenID)
}

func UpgradeContextWithUser(ctx context.Context, user types.UserDetails) context.Context {
	return context.WithValue(ctx, ucKey, user)
}
