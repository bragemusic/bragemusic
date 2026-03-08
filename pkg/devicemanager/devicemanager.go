package devicemanager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/sse"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type DeviceManager struct {
	sseDispatch sse.Dispatcher
	db          database.DeviceFace
}

func (d DeviceManager) ListActiveDevices(ctx context.Context, userID uuid.UUID) ([]types.Device, error) {
	d.sseDispatch.Broadcast(sse.Event{
		ID:   userID,
		Type: "kalaskula broadcast",
		Data: map[string]string{"kaka": "gott"},
	})

	d.sseDispatch.SendToDevice(uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111112")), sse.Event{
		ID:   userID,
		Type: "kalaskula till 2an",
		Data: map[string]string{"kaka": "gott, mkt"},
	})

	devices, err := d.db.ListUsersDevices(ctx, userID)
	if err != nil {
		// FIXME: bragerr
		return nil, err
	}

	for _, adID := range d.sseDispatch.ActiveDevices(userID) {
		for idx := range devices {
			if devices[idx].ID == adID {
				devices[idx].Active = true
				break
			}
		}
	}

	return devices, nil
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

func New(sd sse.Dispatcher, db database.DeviceFace) DeviceManager {
	return DeviceManager{
		sseDispatch: sd,
		db:          db,
	}
}
