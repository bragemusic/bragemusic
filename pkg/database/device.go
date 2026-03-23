package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

type DeviceFace interface {
	executor

	AddDevice(ctx context.Context, device types.Device) error
	DeviceExists(ctx context.Context, deviceID uuid.UUID) (bool, error)
	GetDevice(ctx context.Context, deviceID uuid.UUID) (types.Device, error)
	GetDeviceFromTokenID(ctx context.Context, tokenID uuid.UUID) (types.Device, error)
	GetDeviceFromTokenAndDeviceID(ctx context.Context, deviceID, tokenID uuid.UUID) (device types.Device, err error)
	ListUsersDevices(ctx context.Context, userID uuid.UUID) (devices []types.Device, err error)
	LinkDeviceToken(ctx context.Context, deviceID, tokenID uuid.UUID) error
	UpdateDevice(ctx context.Context, device types.Device) error
	UpdateDeviceLastSeen(ctx context.Context, ipAddress string, deviceID uuid.UUID) error

	//
}

func (d Database) AddDevice(ctx context.Context, device types.Device) error {
	if device.ID == uuid.Nil {
		return errors.New("no device id provided")
	}

	now := time.Now()
	device.CreatedAt = now
	device.UpdatedAt = now

	query := `
        INSERT INTO devices (
            id, name, type, interface, icon, user_id, supports_playback, platform, version, last_ip, last_seen,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		device.ID,
		device.Name,
		device.Type,
		device.Interface,
		device.Icon,
		device.UserID,
		device.SupportsPlayback,
		device.Platform,
		device.Version,
		device.LastIP,
		device.LastSeen,
		device.CreatedAt,
		device.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) DeviceExists(ctx context.Context, deviceID uuid.UUID) (bool, error) {
	if deviceID == uuid.Nil {
		return false, errors.New("no device id provided")
	}

	query := `
		SELECT 1
		FROM devices
		WHERE id = ?
		LIMIT 1;
	`

	var exists int
	err := sqlx.GetContext(ctx, d.ext, &exists, query, deviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d Database) GetDevice(ctx context.Context, deviceID uuid.UUID) (types.Device, error) {
	query := `
		SELECT *
		FROM devices
		WHERE id = ?
		LIMIT 1;
	`

	var device types.Device
	err := sqlx.GetContext(ctx, d.ext, &device, query, deviceID)
	if err != nil {
		return types.Device{}, err
	}

	return device, nil
}

func (d Database) GetDeviceFromTokenAndDeviceID(ctx context.Context, deviceID, tokenID uuid.UUID) (types.Device, error) {
	query := `
		SELECT d.*
		FROM devices d
		JOIN device_tokens dt ON dt.device_id = d.id
		WHERE d.id = ?
		AND dt.token_id = ?
		LIMIT 1;
	`

	var device types.Device
	err := sqlx.GetContext(ctx, d.ext, &device, query, deviceID, tokenID)
	if err != nil {
		return types.Device{}, err
	}

	return device, nil
}

func (d Database) GetDeviceFromTokenID(ctx context.Context, tokenID uuid.UUID) (types.Device, error) {
	query := `
		SELECT d.*
		FROM devices d
		JOIN device_tokens dt ON dt.device_id = d.id
		WHERE dt.token_id = ?
		LIMIT 1;
	`

	var device types.Device
	err := sqlx.GetContext(ctx, d.ext, &device, query, tokenID)
	if err != nil {
		return types.Device{}, err
	}

	return device, nil
}

func (d Database) LinkDeviceToken(ctx context.Context, deviceID, tokenID uuid.UUID) error {
	if deviceID == uuid.Nil {
		return errors.New("no device id provided")
	}

	if tokenID == uuid.Nil {
		return errors.New("no token id provided")
	}

	query := `
		INSERT INTO device_tokens (device_id, token_id, created_at)
		VALUES (?, ?, ?)
		ON CONFLICT(device_id, token_id) DO NOTHING;
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		deviceID,
		tokenID,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) ListUsersDevices(ctx context.Context, userID uuid.UUID) (devices []types.Device, err error) {
	query := `
        SELECT *
        FROM devices
        WHERE user_id = ?
    `

	err = sqlx.SelectContext(ctx, d.ext, &devices, query, userID)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (d Database) UpdateDevice(ctx context.Context, device types.Device) error {
	if device.ID == uuid.Nil {
		return errors.New("no device id provided")
	}

	device.UpdatedAt = time.Now()

	query := `
        UPDATE devices SET
            name = :name,
            type = :type,
            interface = :interface,
            icon = :icon,
            supports_playback = :supports_playback,
            platform = :platform,
            version = :version,
            last_ip = :last_ip,
            last_seen = :last_seen,
            updated_at = :updated_at
        WHERE id = :id;
    `

	_, err := sqlx.NamedExecContext(ctx, d.ext, query, device)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) UpdateDeviceLastSeen(ctx context.Context, ipAddress string, deviceID uuid.UUID) error {
	updatedAt := time.Now()

	query := `
        UPDATE devices SET
            last_ip = ?,
            last_seen = ?,
            updated_at = ?
        WHERE id = ?;
    `

	_, err := d.ext.ExecContext(ctx, query, ipAddress, updatedAt, updatedAt, deviceID)
	if err != nil {
		return err
	}

	return nil
}
