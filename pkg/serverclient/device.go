package serverclient

import (
	"context"
	"net/url"

	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) DevicePlayPause(ctx context.Context, deviceID uuid.UUID) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "play-pause")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, nil, nil); err != nil {
		return err
	}

	return nil
}
