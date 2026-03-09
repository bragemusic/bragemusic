package device

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/sse"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func idFilePath(userID uuid.UUID) (string, error) {
	return xdg.StateFile(filepath.Join("brage", "users", userID.String(), "device_id"))
}

type DeviceAgent struct {
	sc         *serverclient.ServerClient
	deviceMeta types.DeviceBase
	userID     uuid.UUID

	typeRecievers     map[types.SSEventType][]sse.EventHandler
	categoryRecievers map[string][]sse.EventHandler

	log  *slog.Logger
	berr bragerr.BragErrFactory
}

func (a *DeviceAgent) SubscribeToEventTypes(handler sse.EventHandler, eventType ...types.SSEventType) {
	for _, et := range eventType {
		a.typeRecievers[et] = append(a.typeRecievers[et], handler)
	}
}

func (a *DeviceAgent) SubscribeToEventCategory(handler sse.EventHandler, eventCategory string) {
	a.categoryRecievers[eventCategory] = append(a.categoryRecievers[eventCategory], handler)
}

func (a *DeviceAgent) loadLocalDeviceID(userID uuid.UUID) (*uuid.UUID, error) {
	path, err := idFilePath(userID)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	uid, err := uuid.FromString(string(b))
	if err != nil {
		return nil, err
	}

	return &uid, nil
}

func (a *DeviceAgent) saveLocalDeviceID(userID, deviceID uuid.UUID) error {
	path, err := idFilePath(userID)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(deviceID.String()), 0o600); err != nil {
		return err
	}

	return nil
}

func (a *DeviceAgent) handleEvent(ctx context.Context, e types.SSEvent[any]) {
	a.log.DebugContext(ctx, "event recieved", "type", e.Type, "id", e.ID, "data", e.Data)

	handlers, ok := a.typeRecievers[e.Type]
	if ok {
		for _, h := range handlers {
			h(ctx, e)
		}
	}

	category := strings.Split(string(e.Type), ".")[0]
	catHandlers, ok := a.categoryRecievers[category]
	if ok {
		for _, h := range catHandlers {
			h(ctx, e)
		}
	}
}

func (a *DeviceAgent) SubscribeDeviceEvents(ctx context.Context) error {
	deviceID, err := a.loadLocalDeviceID(a.userID)
	if err != nil {
		return err
	}

	newDeviceID, err := a.sc.RegisterDevice(ctx, deviceID, a.deviceMeta)
	if err != nil {
		return err
	}

	if err := a.saveLocalDeviceID(a.userID, newDeviceID); err != nil {
		return err
	}

	if err := a.sc.SubscribeDeviceEvents(ctx, newDeviceID, a.handleEvent); err != nil {
		return err
	}

	return nil
}

func NewAgent(slogHandler slog.Handler, sc *serverclient.ServerClient, userID uuid.UUID, meta types.DeviceBase) DeviceAgent {
	return DeviceAgent{
		log:               slog.New(slogHandler).With("service", "device.agent"),
		berr:              bragerr.NewFactory("device.agent"),
		sc:                sc,
		deviceMeta:        meta,
		userID:            userID,
		typeRecievers:     map[types.SSEventType][]sse.EventHandler{},
		categoryRecievers: map[string][]sse.EventHandler{},
	}
}
