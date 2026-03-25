package device

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
)

type (
	callingDeviceIDContextKeyType struct{}
)

var cdKey = callingDeviceIDContextKeyType{}

func UpgradeContextWithCallingDeviceID(ctx context.Context, deviceID uuid.UUID) context.Context {
	return context.WithValue(ctx, cdKey, deviceID)
}

func CallingDeviceIDFromContext(ctx context.Context) (uuid.UUID, error) {
	td, ok := ctx.Value(cdKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("no calling device token in context")
	}

	return td, nil
}
