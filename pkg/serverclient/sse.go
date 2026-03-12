package serverclient

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bragemusic/core/pkg/sse"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *ServerClient) RegisterDevice(ctx context.Context, deviceID *uuid.UUID, deviceDetails types.DeviceBase) (newDeviceID uuid.UUID, err error) {
	u, err := url.JoinPath(s.baseUrl, DEVICES_PATH...)
	if err != nil {
		return uuid.Nil, err
	}

	req := types.ReqDevicesRegister{
		ID:         deviceID,
		DeviceBase: deviceDetails,
	}

	resp := types.RespDevicesRegister{}
	if err := s.doPostJson(ctx, u, req, &resp); err != nil {
		return uuid.Nil, err
	}

	s.deviceID = &resp.DeviceID

	return resp.DeviceID, nil
}

func (s ServerClient) SubscribeDeviceEvents(ctx context.Context, deviceID uuid.UUID, handler sse.EventHandler) error {
	go func() {
		backoff := time.Second

		for {
			err := s.consumeSSE(ctx, deviceID, handler)
			if err != nil && ctx.Err() == nil {
				s.log.ErrorContext(ctx, "server events connection failed. Retrying.", "error", err.Error(), "backoff", backoff)
				time.Sleep(backoff)

				if backoff < 30*time.Second {
					backoff *= 2
				}

				continue
			}
		}
	}()

	return nil
}

func (s ServerClient) consumeSSE(ctx context.Context, deviceID uuid.UUID, handler sse.EventHandler) error {
	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, deviceID.String(), "events")...)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	reader := bufio.NewScanner(resp.Body)

	for reader.Scan() {
		line := reader.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			raw := strings.TrimSpace(line[5:])
			var evt types.SSEvent[any]
			if err := json.Unmarshal([]byte(raw), &evt); err != nil {
				s.log.ErrorContext(ctx, "Failed to unmarshal event", "error", err.Error())
				continue
			}
			handler(ctx, evt)
		}
	}

	return reader.Err()
}

func (s ServerClient) UpdatePlayContext(ctx context.Context, pc types.PlayContext) error {
	if s.deviceID == nil {
		return errors.New("device not registered")
	}

	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, s.deviceID.String(), "playcontext")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, pc.DTO(*s.deviceID), nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) UpdatePlaybackState(ctx context.Context, ps types.PlaybackState) error {
	if s.deviceID == nil {
		return errors.New("device not registered")
	}

	u, err := url.JoinPath(s.baseUrl, append(DEVICES_PATH, s.deviceID.String(), "playbackstate")...)
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, ps.DTO(*s.deviceID), nil); err != nil {
		return err
	}

	return nil
}
