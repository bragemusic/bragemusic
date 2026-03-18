package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) deviceRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{deviceID}/events", s.sseHub.EventsHandler(), nil, routes.RouteMeta{
			Summary:             "Subscribe to SSE events.",
			Description:         "Streams events from the server to the client.",
			ExpectedDescription: "Streamed event data",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),

		routes.New("GET", "/", s.listDevices(), nil, routes.RouteMeta{
			Summary:             "List user's devices.",
			Description:         "List user's devices. Can filter for only active or all.",
			ExpectedDescription: "Device data",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),

		routes.New("POST", "/", s.registerDevice(), nil, routes.RouteMeta{
			Summary:             "Register or update a device.",
			Description:         "Registers a device if not existing, otherwise updates the exisiting one.",
			ExpectedDescription: "Succesfull registration",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusCreated,
		}),

		routes.New("POST", "/{deviceID}/player/play-pause", s.devicePlayerPlayPause(), nil, routes.RouteMeta{
			Summary:             "Sends a play/pause command to player.",
			Description:         "Sends a command to tell the selected device to play or pause. If the device does not support playback error is returned.",
			ExpectedDescription: "Command sent",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),

		routes.New("POST", "/{deviceID}/player/state", s.devicePlayerSetState(), nil, routes.RouteMeta{
			Summary:             "Sends a playstate command to player.",
			Description:         "Sends a command to tell the selected device to start a new playstate. If the device does not support playback error is returned.",
			ExpectedDescription: "Command sent",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),

		routes.New("POST", "/{deviceID}/playcontext", s.deviceUpdatePlayerPlayContext(), nil, routes.RouteMeta{
			Summary:             "Update the device current play context.",
			Description:         "Update the playcontext for the selected device. Can only be accessed if the device ID is belonging to the token the user is logged in with.",
			ExpectedDescription: "PlayContext updated",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),

		routes.New("POST", "/{deviceID}/playbackstate", s.deviceUpdatePlayerPlaybackState(), nil, routes.RouteMeta{
			Summary:             "Update the device current playback state.",
			Description:         "Update the playback statefor the selected device. Can only be accessed if the device ID is belonging to the token the user is logged in with.",
			ExpectedDescription: "PlaybackState updated",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

// FIXME: The routes object should probably be divided in to 3 parts:
// - listdevices, registerdevice (standard auth)
// - play-pause, start track, next track and so on: deviceID is owned by the user
// - update playcontext, playback state: deviceID is the one that correspondes with the used token. This might require unique on the device - token relationship.

func (s *Server) listDevices() routes.RouteFunc[ReqNoContent, types.ListPayload[types.DeviceDetailed]] {
	return func(ctx context.Context, req ReqNoContent, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.DeviceDetailed]], err error) {
		d, err := s.devicemgr.ListActiveDevices(ctx, user.ID)
		if err != nil {
			return types.Response[types.ListPayload[types.DeviceDetailed]]{}, err
		}

		return types.Response[types.ListPayload[types.DeviceDetailed]]{
			Payload: types.ListPayload[types.DeviceDetailed]{
				Items: d,
				Count: len(d),
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) registerDevice() routes.RouteFunc[ReqDevicesRegister, types.RespDevicesRegister] {
	return func(ctx context.Context, req ReqDevicesRegister, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.RespDevicesRegister], err error) {
		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			return types.Response[types.RespDevicesRegister]{}, err
		}

		device := types.Device{
			DeviceBase: req.DeviceBase,
			UserID:     user.ID,
			LastIP:     "must.add.ip.here",
		}

		id, err := s.devicemgr.RegisterOrUpdateDevice(ctx, req.ID, tokenID, user.ID, device)
		if err != nil {
			return types.Response[types.RespDevicesRegister]{}, err
		}

		if req.ID == nil {
			return types.Response[types.RespDevicesRegister]{
				Payload: types.RespDevicesRegister{
					DeviceID: id,
				},
				Status: http.StatusCreated,
			}, nil
		} else {
			return types.Response[types.RespDevicesRegister]{
				Payload: types.RespDevicesRegister{
					DeviceID: id,
				},
				Status: http.StatusOK,
			}, nil
		}
	}
}

func (s *Server) devicePlayerPlayPause() routes.RouteFunc[ReqDevicesGet, types.NoResponse] {
	return func(ctx context.Context, req ReqDevicesGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		// FIXME
		callingDeviceID := uuid.Nil
		if err := s.devicemgr.PlayerPlayPause(ctx, req.DeviceID, callingDeviceID, user.ID); err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) deviceUpdatePlayerPlayContext() routes.RouteFunc[ReqDevicesUpdatePlayContext, types.NoResponse] {
	return func(ctx context.Context, req ReqDevicesUpdatePlayContext, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.devicemgr.UpdatePlayerPlayContext(ctx, req.PlayContextDTO, req.DeviceID, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) deviceUpdatePlayerPlaybackState() routes.RouteFunc[ReqDevicesUpdatePlaybackState, types.NoResponse] {
	return func(ctx context.Context, req ReqDevicesUpdatePlaybackState, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.devicemgr.UpdatePlayerPlaybackState(ctx, req.PlaybackStateDTO, req.DeviceID, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) devicePlayerSetState() routes.RouteFunc[ReqDevicesPlayerSetState, types.NoResponse] {
	return func(ctx context.Context, req ReqDevicesPlayerSetState, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		// FIXME
		callingDeviceID := uuid.Nil
		if err := s.devicemgr.PlayerSetState(ctx, req.PlayerState, req.DeviceID, callingDeviceID, user.ID); err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusOK,
		}, nil
	}
}
