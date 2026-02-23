package client

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

func (c *ClientStreaming) AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error {
	track, err := c.GetTrackDetailed(ctx, trackID, albumID)
	if err != nil {
		return err
	}

	c.AudioPlayer.AddTrackToQueue(ctx, track)
	return nil
}
