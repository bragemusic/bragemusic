package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) DeviceNextTrack(ctx context.Context, deviceID uuid.UUID) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "next")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, nil, nil); err != nil {
		return err
	}

	return nil
}

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

func (s ServerClient) DevicePreviousTrack(ctx context.Context, deviceID uuid.UUID) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "previous")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, nil, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DevicePlayerSetRepeat(ctx context.Context, deviceID uuid.UUID, repeatType types.RepeatType) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "repeat")...)
	if err != nil {
		return err
	}

	req := types.SetDevicePlayerRepeatReq{
		Type: repeatType,
	}

	if err := s.doPostJson(ctx, u, req, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DeviceSetPlayerState(ctx context.Context, deviceID uuid.UUID, pc types.PlayerState) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "state")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, pc, nil); err != nil {
		return err
	}

	return nil
}
