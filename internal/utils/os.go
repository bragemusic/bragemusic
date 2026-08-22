package utils

import (
	"context"
	"os"

	"github.com/bragemusic/bragemusic/pkg/config"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func SetupFolderStructure(ctx context.Context, cfg config.ClientConfig) error {
	if cfg.General.ClientType == types.DeviceTypeStreaming {
		return nil
	}

	if err := os.MkdirAll(cfg.Paths.ConfigDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.Paths.ImageDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.Paths.MusicDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}
