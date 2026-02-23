package syncer

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
)

type NoSync struct{}

func (c *NoSync) RegisterSyncInProgressCallback(f func(bool)) {
}

func (c NoSync) SupportsSync() bool {
	return false
}

func (c NoSync) Sync(ctx context.Context, userID uuid.UUID) error {
	return errors.New("sync not available")
}
