package device

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/sse"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

type DeviceManager struct {
	sseDispatch sse.Dispatcher
	db          database.DeviceFace

	log *slog.Logger

	playerStates map[uuid.UUID]types.PlayerStateDTO
}

// FIXME
func (d DeviceManager) hasAccess(ctx context.Context, deviceID, userID uuid.UUID) (bool, error) {
	return true, nil
}

// FIXME
func (d DeviceManager) isConnected(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) (bool, error) {
	return true, nil
}

func (d DeviceManager) PlayerPlayPause(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	hasAccess, err := d.hasAccess(ctx, targetDeviceID, userID)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.New("FIXME no access")
	}

	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerPlayPause()); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) UpdatePlayerPlayContext(ctx context.Context, pc types.PlayContextDTO, deviceID, userID uuid.UUID) error {
	// FIXME: security

	state, found := d.playerStates[deviceID]
	if !found {
		state = types.PlayerStateDTO{}
	}

	state.Context = pc

	d.playerStates[deviceID] = state

	err := d.sseDispatch.SendToUser(userID, types.SSEPlayerPlayContext(pc))
	if err != nil {
		d.log.WarnContext(ctx, "could not send play context to user's devices", "userID", userID.String(), "deviceID", deviceID.String(), "error", err.Error())
	}

	return nil
}

func (d DeviceManager) UpdatePlayerPlaybackState(ctx context.Context, ps types.PlaybackStateDTO, deviceID, userID uuid.UUID) error {
	// FIXME: security

	state, found := d.playerStates[deviceID]
	if !found {
		state = types.PlayerStateDTO{}
	}

	state.Playback = ps

	d.playerStates[deviceID] = state

	err := d.sseDispatch.SendToUser(userID, types.SSEPlayerPlaybackState(ps))
	if err != nil {
		d.log.WarnContext(ctx, "could not send playback state to user's devices", "userID", userID.String(), "deviceID", deviceID.String(), "error", err.Error())
	}

	return nil
}

func (d DeviceManager) ListActiveDevices(ctx context.Context, userID uuid.UUID) ([]types.DeviceDetailed, error) {
	// d.sseDispatch.Broadcast(sse.Event{
	// 	ID:   userID,
	// 	Type: "kalaskula broadcast",
	// 	Data: map[string]string{"kaka": "gott"},
	// })

	// d.sseDispatch.SendToDevice(uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111112")), sse.Event{
	// 	ID:   userID,
	// 	Type: "kalaskula till 2an",
	// 	Data: map[string]string{"kaka": "gott, mkt"},
	// })

	devices, err := d.db.ListUsersDevices(ctx, userID)
	if err != nil {
		// FIXME: bragerr
		return nil, err
	}

	devicesDetailed := lo.Map(devices, func(item types.Device, index int) types.DeviceDetailed {
		return types.DeviceDetailed{Device: item}
	})

	// NOTE: Should probably add playstater here somwhoe. Do a DeviceDetailed type

	for _, adID := range d.sseDispatch.ActiveDevices(userID) {
		for idx := range devicesDetailed {
			if devicesDetailed[idx].ID == adID {
				devicesDetailed[idx].Active = true
				ps, ok := d.playerStates[devicesDetailed[idx].ID]
				if ok {
					devicesDetailed[idx].PlayerState = utils.Ptr(ps)
				}
				break
			}
		}
	}

	return devicesDetailed, nil
}

func (d DeviceManager) RegisterOrUpdateDevice(ctx context.Context, id *uuid.UUID, tokenID, userID uuid.UUID, device types.Device) (deviceID uuid.UUID, err error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer tx.Rollback()

	fmt.Println("ueue", id)
	if id != nil {
		existing, err := tx.GetDevice(ctx, *id)
		if err != nil {
			// FIXME: bragerr
			return uuid.Nil, err
		}

		if existing.UserID != userID {
			// FIXME: bragerr
			return uuid.Nil, errors.New("not your device. Fast inte sa avslojande")
		}

		existing.Name = device.Name
		existing.Type = device.Type
		existing.Interface = device.Interface
		existing.SupportsPlayback = device.SupportsPlayback
		existing.Platform = device.Platform
		existing.Version = device.Version
		existing.LastIP = device.LastIP
		existing.UpdatedAt = time.Now()

		fmt.Println("uppdatera", device.Name, existing.Name)
		if err = tx.UpdateDevice(ctx, existing); err != nil {
			// FIXME: bragerr
			return uuid.Nil, err
		}

		if err = tx.LinkDeviceToken(ctx, *id, tokenID); err != nil {
			// FIXME: bragerr
			return uuid.Nil, err
		}

		if err = tx.Commit(); err != nil {
			return uuid.Nil, err
		}

		return *id, nil
	}

	newID, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	device.ID = newID
	device.UserID = userID

	if err = tx.AddDevice(ctx, device); err != nil {
		// FIXME: bragerr
		return uuid.Nil, err
	}

	if err = tx.LinkDeviceToken(ctx, newID, tokenID); err != nil {
		// FIXME: bragerr
		return uuid.Nil, err
	}

	if err = tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

func NewManager(sd sse.Dispatcher, db database.DeviceFace, slogHandler slog.Handler) DeviceManager {
	return DeviceManager{
		sseDispatch:  sd,
		db:           db,
		log:          slog.New(slogHandler).With("service", "device.manager"),
		playerStates: map[uuid.UUID]types.PlayerStateDTO{},
	}
}
