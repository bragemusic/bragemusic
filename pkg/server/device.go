package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
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
	}
}

func (s *Server) listDevices() routes.RouteFunc[ReqNoContent, types.ListPayload[types.Device]] {
	return func(ctx context.Context, req ReqNoContent, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.Device]], err error) {
		d, err := s.devicemgr.ListActiveDevices(ctx, user.ID)
		if err != nil {
			return types.Response[types.ListPayload[types.Device]]{}, err
		}

		return types.Response[types.ListPayload[types.Device]]{
			Payload: types.ListPayload[types.Device]{
				Items: d,
				Count: len(d),
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) registerDevice() routes.RouteFunc[ReqDevicesRegister, RespDevicesRegister] {
	return func(ctx context.Context, req ReqDevicesRegister, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[RespDevicesRegister], err error) {
		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			return types.Response[RespDevicesRegister]{}, err
		}

		device := types.Device{
			Name:             req.Name,
			Type:             req.Type,
			Interface:        req.Interface,
			UserID:           user.ID,
			SupportsPlayback: req.SupportsPlayback,
			Platform:         req.Platform,
			Version:          req.Version,
			LastIP:           "must.add.ip.here",
		}

		fmt.Println("kalas!", req)
		id, err := s.devicemgr.RegisterOrUpdateDevice(ctx, req.DeviceID, tokenID, user.ID, device)
		if err != nil {
			return types.Response[RespDevicesRegister]{}, err
		}

		if req.DeviceID == nil {
			return types.Response[RespDevicesRegister]{
				Payload: RespDevicesRegister{
					DeviceID: id,
				},
				Status: http.StatusCreated,
			}, nil
		} else {
			return types.Response[RespDevicesRegister]{
				Payload: RespDevicesRegister{
					DeviceID: id,
				},
				Status: http.StatusOK,
			}, nil
		}
	}
}
