package app

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
)

func (a *App) ListDevices() []types.DeviceDetailed {
	if a.client == nil {
		return []types.DeviceDetailed{}
	}

	devices, err := a.client.ListDevices(a.ctx)
	if err != nil {
		a.handleError(err)
		return []types.DeviceDetailed{}
	}

	if len(devices) == 0 {
		return []types.DeviceDetailed{}
	}

	return devices
}

// func (a *App) ListActivePlayableDevices() []types.DeviceDetailed {
// 	// FIXME: Should not be here
// 	devices, err := a.client.ListDevices(a.ctx)
// 	if err != nil {
// 		a.handleError(err)
// 		return []types.DeviceDetailed{}
// 	}

// 	if len(devices) == 0 {
// 		return []types.DeviceDetailed{}
// 	}

// 	fDevices := []types.DeviceDetailed{}

// 	for _, d := range devices {
// 		if d.Active && d.SupportsPlayback {
// 			fDevices = append(fDevices, d)
// 		}
// 	}

// 	if len(fDevices) == 0 {
// 		return []types.DeviceDetailed{}
// 	}

// 	return fDevices
// }

func (a *App) ConnectDevice(id string) {
	uID, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.ConnectDevice(a.ctx, uID)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) DisconnectDevice() {
	err := a.client.DisconnectDevice(a.ctx)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) GetConnectedDeviceID() *string {
	id := a.client.GetConnectedDevice()
	if id == nil {
		return nil
	}
	return utils.Ptr(id.String())
}

func (a *App) RemoveDeviceToken(deviceID string) {
	uid, err := uuid.FromString(deviceID)
	if err != nil {
		a.handleError(err)
		return
	}
	_ = uid

	err = a.client.DeleteDeviceToken(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) RemoveDevice(deviceID string) {
	uid, err := uuid.FromString(deviceID)
	if err != nil {
		a.handleError(err)
		return
	}
	_ = uid

	err = a.client.DeleteDevice(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) RemoveDeviceAndToken(deviceID string) {
	uid, err := uuid.FromString(deviceID)
	if err != nil {
		a.handleError(err)
		return
	}
	_ = uid

	err = a.client.DeleteDeviceAndToken(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}
