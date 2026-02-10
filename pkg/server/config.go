package server

import (
	"errors"
	"fmt"

	"github.com/bragemusic/core/pkg/utils"
)

type Paths struct {
	ConfigDir string `toml:"config_dir" desc:"Dir where server files are stored."`
	ImageDir  string `toml:"image_dir" desc:"Dir where image assets are stored."`
	MusicDir  string `toml:"music_dir" desc:"Dir where music files are stored."`
	ImportDir string `toml:"import_dir" desc:"Dir where imported albums and tracks will be saved before processing."`
}

type Admin struct {
	Email    string `toml:"email" desc:"Default user, with admin rights, email. Defaults to 'admin@example.com'"`
	Username string `toml:"username" desc:"Default user, with admin rights, username. Defaults to 'admin'"`
	Password string `toml:"password" desc:"Default user, with admin rights, password. Defaults to 'password'"`
}

type Config struct {
	Admin Admin `toml:"admin"`
	Paths Paths `toml:"paths"`
	Port  int   `toml:"port" desc:"Port of the server. Defaults to 3000."`
}

var defaultConfig = Config{
	Admin: Admin{
		Email:    "admin@example.com",
		Username: "admin",
		Password: "password",
	},
	Paths: Paths{},
	Port:  3000,
}

func verify(cfg Config) error {
	errs := []error{}

	if cfg.Paths.ImageDir == "" {
		errs = append(errs, errors.New("Paths.ImageDir not set"))
	}

	if cfg.Paths.MusicDir == "" {
		errs = append(errs, errors.New("Paths.MusicDir not set"))
	}

	if cfg.Paths.ConfigDir == "" {
		errs = append(errs, errors.New("Paths.ConfigDir not set"))
	}

	if cfg.Paths.ImportDir == "" {
		errs = append(errs, errors.New("Paths.ImportDir not set"))
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
