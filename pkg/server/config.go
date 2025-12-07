package server

import (
	"errors"
	"fmt"

	"github.com/bragemusic/core/pkg/utils"
)

type Paths struct {
	ConfigDir string `toml:"config_dir"`
	ImageDir  string `toml:"image_dir"`
	MusicDir  string `toml:"music_dir"`
}

type Config struct {
	Paths Paths `toml:"paths"`
	Port  int   `toml:"port"`
}

var defaultConfig = Config{
	Paths: Paths{},
	Port:  3000,
}

func verify(cfg Config) error {
	errs := []error{}

	if cfg.Paths.ImageDir == "" {
		errs = append(errs, errors.New("Paths.ImageDir not set"))
	}

	if cfg.Paths.ConfigDir == "" {
		errs = append(errs, errors.New("Paths.ConfigDir not set"))
	}

	if cfg.Port < 1024 || cfg.Port > 49151 {
		errs = append(errs, fmt.Errorf("port %d not allowed. Use in range 1024-49151", cfg.Port))
	}

	// TODO: Add more verification
	return errors.Join(errs...)
}

func GetConfig() (Config, error) {
	c := defaultConfig

	err := utils.AddFromEnv(&c)
	if err != nil {
		return Config{}, err
	}

	if err = verify(c); err != nil {
		return Config{}, err
	}

	return c, nil
}
