package auth

import (
	"context"
	"errors"

	"github.com/bragemusic/core/pkg/types"
)

type userContextKeyType struct{}

var ucKey = userContextKeyType{}

func UserFromContext(ctx context.Context) (types.UserDetails, error) {
	ud, ok := ctx.Value(ucKey).(types.UserDetails)
	if !ok {
		return types.UserDetails{}, errors.New("no user in context")
	}

	return ud, nil
}

func UpgradeContextWithUser(ctx context.Context, user types.UserDetails) context.Context {
	return context.WithValue(ctx, ucKey, user)
}
