package device

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/bragerr"
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

	log  *slog.Logger
	berr bragerr.BragErrFactory

	playerStates map[uuid.UUID]types.PlayerStateDTO
}

func (d DeviceManager) MiddlewareUserAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		deviceID, err := utils.GetURLParameter[uuid.UUID](ctx, "deviceID")
		if err != nil {
			bragerr.HandleHttpResponse(ctx, d.berr.ParamWrongFormat(err, "deviceID", "uuid"), w, d.log)
			return
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
			return
		}

		hasAccess, err := d.hasAccess(ctx, deviceID, user.ID)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
			return
		}

		if !hasAccess {
			bragerr.HandleHttpResponse(ctx, d.berr.Unauthenticated(errors.New("user does not have access to device")), w, d.log)
			return
		}

		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
		}

		device, err := d.db.GetDeviceFromTokenID(ctx, tokenID)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
		}

		ctx = UpgradeContextWithCallingDeviceID(ctx, device.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (d DeviceManager) MiddlewareDeviceAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		deviceID, err := utils.GetURLParameter[uuid.UUID](ctx, "deviceID")
		if err != nil {
			bragerr.HandleHttpResponse(ctx, d.berr.ParamWrongFormat(err, "deviceID", "uuid"), w, d.log)
			return
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
			return
		}

		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
		}

		hasAccess, err := d.hasDeviceAccess(ctx, deviceID, user.ID, tokenID)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, d.log)
			return
		}

		if !hasAccess {
			bragerr.HandleHttpResponse(ctx, d.berr.Unauthenticated(errors.New("device does not have access to device")), w, d.log)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (d DeviceManager) hasAccess(ctx context.Context, deviceID, userID uuid.UUID) (bool, error) {
	device, err := d.db.GetDevice(ctx, deviceID)
	if err != nil {
		return false, d.berr.DatabaseError(err, types.EntityDevice, &deviceID)
	}

	if device.UserID != userID {
		return false, nil
	}

	return true, nil
}

func (d DeviceManager) hasDeviceAccess(ctx context.Context, deviceID, userID, tokenID uuid.UUID) (bool, error) {
	device, err := d.db.GetDeviceFromTokenAndDeviceID(ctx, deviceID, tokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, d.berr.DatabaseError(err, types.EntityDevice, &deviceID)
	}

	if deviceID != device.ID || userID != device.UserID {
		return false, nil
	}

	return true, nil
}

// FIXME
func (d DeviceManager) isConnected(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) (bool, error) {
	return true, nil
}

func (d DeviceManager) PlayerNextTrack(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerNextTrack()); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) PlayerPlayPause(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
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

func (d DeviceManager) PlayerPreviousTrack(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerPreviousTrack()); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) PlayerSetRepeat(ctx context.Context, repeatType types.RepeatType, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerSetRepeat(repeatType)); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) PlayerSetShuffle(ctx context.Context, active bool, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerSetShuffle(active)); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) PlayerSetState(ctx context.Context, ps types.PlayerState, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerSetState(ps)); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) PlayerStop(ctx context.Context, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerStop()); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) PlayerAddToQueue(ctx context.Context, track types.TrackDetailed, targetDeviceID, callingDeviceID, userID uuid.UUID) error {
	isConnected, err := d.isConnected(ctx, targetDeviceID, callingDeviceID, userID)
	if err != nil {
		return err
	}

	if !isConnected {
		return errors.New("FIXME not connected")
	}

	if err = d.sseDispatch.SendToDevice(targetDeviceID, types.SSEPlayerAddToQueue(track)); err != nil {
		return err
	}

	return nil
}

func (d DeviceManager) UpdatePlayerPlayContext(ctx context.Context, pc types.PlayContextDTO, deviceID, userID uuid.UUID) error {
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
	devices, err := d.db.ListUsersDevices(ctx, userID)
	if err != nil {
		// FIXME: bragerr
		return nil, err
	}

	devicesDetailed := lo.Map(devices, func(item types.Device, index int) types.DeviceDetailed {
		return types.DeviceDetailed{Device: item}
	})

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
		existing.Icon = device.Icon
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

func (d DeviceManager) GetDeviceToken(ctx context.Context, deviceID, userID uuid.UUID) (types.DeviceToken, error) {
	device, err := d.db.GetDevice(ctx, deviceID)
	if err != nil {
		return types.DeviceToken{}, err
	}

	if device.UserID != userID {
		return types.DeviceToken{}, d.berr.ItemAccessDenied(errors.New("user is not the owner of the device"), types.EntityDevice, deviceID)
	}

	dt, err := d.db.GetDeviceTokenFromDeviceID(ctx, deviceID)
	if err != nil {
		return types.DeviceToken{}, d.berr.DatabaseError(err, types.EntityDeviceToken, nil)
	}

	return dt, nil
}

func NewManager(sd sse.Dispatcher, db database.DeviceFace, slogHandler slog.Handler) DeviceManager {
	return DeviceManager{
		sseDispatch:  sd,
		db:           db,
		log:          slog.New(slogHandler).With("service", "device.manager"),
		berr:         bragerr.NewFactory("device.manager"),
		playerStates: map[uuid.UUID]types.PlayerStateDTO{},
	}
}
