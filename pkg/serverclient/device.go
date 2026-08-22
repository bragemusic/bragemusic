package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/bragemusic/pkg/types"
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

func (s ServerClient) DevicePlayerAddToQueue(ctx context.Context, deviceID uuid.UUID, track types.TrackDetailed) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "queue")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, track, nil); err != nil {
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

func (s ServerClient) DevicePlayerSetShuffle(ctx context.Context, deviceID uuid.UUID, active bool) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "shuffle")...)
	if err != nil {
		return err
	}

	req := types.SetDevicePlayerShuffleReq{
		Active: active,
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

func (s ServerClient) DevicePlayerStop(ctx context.Context, deviceID uuid.UUID) (err error) {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "player", "stop")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, nil, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DeviceDelete(ctx context.Context, deviceID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String())...)
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DeviceDeleteToken(ctx context.Context, deviceID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "token")...)
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}
