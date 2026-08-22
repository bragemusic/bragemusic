package device

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adrg/xdg"
	"github.com/bragemusic/bragemusic/pkg/bragerr"
	"github.com/bragemusic/bragemusic/pkg/serverclient"
	"github.com/bragemusic/bragemusic/pkg/sse"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type DeviceAgent struct {
	sc         *serverclient.ServerClient
	deviceMeta types.DeviceBase
	userID     uuid.UUID

	devices map[uuid.UUID]types.DeviceDetailed
	mu      sync.RWMutex

	typeRecievers     map[types.SSEventType][]sse.EventHandler
	categoryRecievers map[string][]sse.EventHandler

	clientEventHandlers []sse.EventHandler

	stateFilePath *string
	log           *slog.Logger
	berr          bragerr.BragErrFactory
}

// Get a device (read)
func (a *DeviceAgent) getDevice(id uuid.UUID) (types.DeviceDetailed, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	d, ok := a.devices[id]
	return d, ok
}

// Set or update a device (write)
func (a *DeviceAgent) setDevice(d types.DeviceDetailed) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.devices[d.ID] = d
	a.sendClientEvent(context.Background(), types.SSEDeviceUpdated(a.devices[d.ID].Device))
}

// Delete a device (write)
func (a *DeviceAgent) deleteDevice(id uuid.UUID) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sendClientEvent(context.Background(), types.SSEDeviceDisconnected(a.devices[id].Device))
	delete(a.devices, id)
}

// Replace all devices atomically
func (a *DeviceAgent) replaceDevices(newDevices []types.DeviceDetailed) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.devices = make(map[uuid.UUID]types.DeviceDetailed, len(newDevices))
	for _, d := range newDevices {
		a.devices[d.ID] = d
	}
	// a.sendClientEvent(context.Background(), types.SSEInternalClientAllReplaced(newDevices))
}

func (a *DeviceAgent) idFilePath(userID uuid.UUID) (string, error) {
	if a.stateFilePath != nil {
		return filepath.Join(*a.stateFilePath, "device_id"), nil
	}
	return xdg.StateFile(filepath.Join("brage", "users", userID.String(), "device_id"))
}

func (a *DeviceAgent) sendClientEvent(ctx context.Context, e types.SSEvent) {
	for _, f := range a.clientEventHandlers {
		f(ctx, e)
	}
}

func (a *DeviceAgent) SubscribeToClientEvents(handler sse.EventHandler) {
	a.clientEventHandlers = append(a.clientEventHandlers, handler)
}

func (a *DeviceAgent) SubscribeToEventTypes(handler sse.EventHandler, eventType ...types.SSEventType) {
	for _, et := range eventType {
		a.typeRecievers[et] = append(a.typeRecievers[et], handler)
	}
}

func (a *DeviceAgent) SubscribeToEventCategory(handler sse.EventHandler, eventCategory string) {
	a.categoryRecievers[eventCategory] = append(a.categoryRecievers[eventCategory], handler)
}

func (a *DeviceAgent) handleDeviceEvents(ctx context.Context, e types.SSEvent) {
	switch e.Type {
	case types.SSEventTypeDeviceConnected:
		d, err := types.DecodeEventData[types.Device](e)
		if err != nil {
			a.log.ErrorContext(ctx, "could not decode device data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		a.setDevice(types.DeviceDetailed{
			Device: d,
		})

		a.log.InfoContext(ctx, "new client connection", "device.name", d.Name)
		return

	case types.SSEventTypeDeviceDisconnected:
		d, err := types.DecodeEventData[types.Device](e)
		if err != nil {
			a.log.ErrorContext(ctx, "could not decode device data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		a.deleteDevice(d.ID)

		a.log.InfoContext(ctx, "client disconnecteted", "device.name", d.Name)
		return

	case types.SSEventTypeDeviceDeleted:
		id, err := types.DecodeEventData[uuid.UUID](e)
		if err != nil {
			a.log.ErrorContext(ctx, "could not decode device data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		a.deleteDevice(id)

		a.log.InfoContext(ctx, "client deleted", "device.id", id)
		return

	case types.SSEventTypePlayerPlayContext:
		pc, err := types.DecodeEventData[types.PlayContextDTO](e)
		if err != nil {
			a.log.ErrorContext(ctx, "could not decode playcontext data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		d, ok := a.getDevice(pc.DeviceID)
		if !ok {
			a.log.ErrorContext(ctx, "device not connected, cannot update playcontext", "device.id", pc.DeviceID.String())
			return
		}

		if d.PlayerState == nil {
			d.PlayerState = &types.PlayerStateDTO{}
		}

		d.PlayerState.Context = pc
		a.setDevice(d)
		a.sendClientEvent(ctx, types.SSEDevicePlayContext(pc))

	case types.SSEventTypePlayerPlaybackState:
		ps, err := types.DecodeEventData[types.PlaybackStateDTO](e)
		if err != nil {
			a.log.ErrorContext(ctx, "could not decode playback state data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		d, ok := a.getDevice(ps.DeviceID)
		if !ok {
			a.log.ErrorContext(ctx, "device not connected, cannot update playback state", "device.id", ps.DeviceID.String())
			return
		}

		if d.PlayerState == nil {
			d.PlayerState = &types.PlayerStateDTO{}
		}

		d.PlayerState.Playback = ps
		a.setDevice(d)

		a.sendClientEvent(ctx, types.SSEDevicePlaybackState(ps))
	}
}

func (a *DeviceAgent) loadLocalDeviceID(userID uuid.UUID) (*uuid.UUID, error) {
	path, err := a.idFilePath(userID)
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
	path, err := a.idFilePath(userID)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(deviceID.String()), 0o600); err != nil {
		return err
	}

	return nil
}

func (a *DeviceAgent) handleEvent(ctx context.Context, e types.SSEvent) {
	a.log.DebugContext(ctx, "event recieved", "type", e.Type, "id", e.ID, "data", string(e.Data))

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
		serr, ok := err.(serverclient.ErrStatus)
		if !ok {
			return err
		}

		if serr.Status != http.StatusNotFound {
			return err
		}

		newDeviceID, err = a.sc.RegisterDevice(ctx, nil, a.deviceMeta)
		if err != nil {
			return err
		}

		a.log.InfoContext(ctx, "device does not exist on server. New ID has been given", "old_id", deviceID, "new_id", newDeviceID)
	}

	if err := a.saveLocalDeviceID(a.userID, newDeviceID); err != nil {
		return err
	}

	if err := a.sc.SubscribeDeviceEvents(ctx, newDeviceID, a.handleEvent); err != nil {
		return err
	}

	if _, err := a.ListDevices(ctx); err != nil {
		return err
	}

	return nil
}

func (a *DeviceAgent) ListDevices(ctx context.Context) (devices []types.DeviceDetailed, err error) {
	devices, err = a.sc.ListDevices(ctx)
	if err != nil {
		return nil, err
	}

	a.replaceDevices(devices)

	return devices, nil
}

func (a *DeviceAgent) GetDevice(ctx context.Context, id uuid.UUID) (device types.DeviceDetailed, err error) {
	device, found := a.getDevice(id)
	if !found {
		return types.DeviceDetailed{}, errors.New("device not found")
	}

	return device, nil
}

func (a *DeviceAgent) DeleteDeviceToken(ctx context.Context, deviceID uuid.UUID) error {
	return a.sc.DeviceDeleteToken(ctx, deviceID)
}

func (a *DeviceAgent) DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	return a.sc.DeviceDelete(ctx, deviceID)
}

func (a *DeviceAgent) DeleteDeviceAndToken(ctx context.Context, deviceID uuid.UUID) error {
	if err := a.sc.DeviceDeleteToken(ctx, deviceID); err != nil {
		return err
	}

	return a.sc.DeviceDelete(ctx, deviceID)
}

func NewAgent(slogHandler slog.Handler, sc *serverclient.ServerClient, userID uuid.UUID, meta types.DeviceBase, stateFilePath *string) *DeviceAgent {
	da := &DeviceAgent{
		log:                 slog.New(slogHandler).With("service", "device.agent"),
		berr:                bragerr.NewFactory("device.agent"),
		sc:                  sc,
		deviceMeta:          meta,
		userID:              userID,
		stateFilePath:       stateFilePath,
		typeRecievers:       map[types.SSEventType][]sse.EventHandler{},
		categoryRecievers:   map[string][]sse.EventHandler{},
		devices:             map[uuid.UUID]types.DeviceDetailed{},
		clientEventHandlers: []sse.EventHandler{},
	}

	da.SubscribeToEventCategory(da.handleDeviceEvents, "device")
	da.SubscribeToEventTypes(da.handleDeviceEvents, types.SSEventTypePlayerPlayContext, types.SSEventTypePlayerPlaybackState)

	return da
}
