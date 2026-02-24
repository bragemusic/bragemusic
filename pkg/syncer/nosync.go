package syncer

import (
	"context"
	"errors"
)

type NoSync struct{}

func (c *NoSync) RegisterSyncInProgressCallback(f func(bool)) {
}

func (c NoSync) SupportsSync() bool {
	return false
}

func (c NoSync) Sync(ctx context.Context) error {
	return errors.New("sync not available")
}
