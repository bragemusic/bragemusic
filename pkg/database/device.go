package database

import (
	"context"
	"errors"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

type DeviceFace interface {
	executor

	AddDevice(ctx context.Context, device types.Device) error
	GetDeviceFromTokenID(ctx context.Context, tokenID uuid.UUID) (device types.Device, err error)
	UpdateDevice(ctx context.Context, device types.Device) error

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
            id, name, type, interface, token_id, supports_playback, platform, version, last_ip, last_seen,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err := d.ext.ExecContext(
		ctx,
		query,
		device.ID,
		device.Name,
		device.Type,
		device.Interface,
		device.TokenID,
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

func (d Database) GetDeviceFromTokenID(ctx context.Context, tokenID uuid.UUID) (device types.Device, err error) {
	query := `
        SELECT *
        FROM devices
        WHERE token_id = ?
        LIMIT 1;
    `

	err = sqlx.GetContext(ctx, d.ext, &device, query, tokenID)
	if err != nil {
		return types.Device{}, err
	}

	return device, err
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
