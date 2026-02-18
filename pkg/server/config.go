package server

import (
	"errors"
	"fmt"

	"github.com/adhocore/gronx"
	"github.com/bragemusic/core/pkg/utils"
)

type Paths struct {
	ConfigDir       string `toml:"config_dir" desc:"Dir where server files are stored."`
	ImageDir        string `toml:"image_dir" desc:"Dir where image assets are stored."`
	MusicDir        string `toml:"music_dir" desc:"Dir where music files are stored."`
	ImportDir       string `toml:"import_dir" desc:"Dir where imported albums and tracks will be saved before processing."`
	BackupImportDir string `toml:"backup_import_dir" desc:"Dir where imported albums and tracks will be saved after processing."`
}

type Admin struct {
	Email    string `toml:"email" desc:"Default user, with admin rights, email. Defaults to 'admin@example.com'"`
	Username string `toml:"username" desc:"Default user, with admin rights, username. Defaults to 'admin'"`
	Password string `toml:"password" desc:"Default user, with admin rights, password. Defaults to 'password'"`
}

type AcoustID struct {
	ApiKey string `toml:"api_key" desc:"API key to Acoust ID. Used to identify the files to not rely solely on ID3."`
}

type Wikipedia struct {
	Email string `toml:"email" desc:"Used against the wikipedia API. They require a valid email to make sure you behave."`
}

type Jobs struct {
	Importer   string `toml:"importer" desc:"How often the importer will look for new media files. Cron expression. Defaults to '*/3 * * * *'"`
	MetaSyncer string `toml:"meta_syncer" desc:"How often the meta-syncer will sync the needed metadata. Cron expression. Defaults to '*/3 * * * *'"`
}

type Config struct {
	AcoustID  AcoustID  `toml:"acoust_id"`
	Admin     Admin     `toml:"admin"`
	Jobs      Jobs      `toml:"jobs"`
	Paths     Paths     `toml:"paths"`
	Wikipedia Wikipedia `toml:"wikipedia"`
	Name      string    `toml:"name" desc:"Name of the server. Defaults to 'Brage Music Server'"`
	Port      int       `toml:"port" desc:"Port of the server. Defaults to 3000."`
}

var defaultConfig = Config{
	Admin: Admin{
		Email:    "admin@example.com",
		Username: "admin",
		Password: "password",
	},
	Paths: Paths{},
	Jobs: Jobs{
		Importer:   "*/3 * * * *",
		MetaSyncer: "*/3 * * * *",
	},
	Name: "Brage Music Server",
	Port: 3000,
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

	if cfg.Paths.BackupImportDir == "" {
		errs = append(errs, errors.New("Paths.BackupImportDir not set"))
	}

	if cfg.AcoustID.ApiKey == "" {
		errs = append(errs, errors.New("AcoustID.ApiKey not set"))
	}

	if cfg.Wikipedia.Email == "" {
		errs = append(errs, errors.New("Wikipedia.Email not set"))
	}

	if !gronx.IsValid(cfg.Jobs.Importer) {
		errs = append(errs, errors.New("Jobs.Importer is not a valid Cron Expression"))
	}

	if !gronx.IsValid(cfg.Jobs.MetaSyncer) {
		errs = append(errs, errors.New("Jobs.MetaSyncer is not a valid Cron Expression"))
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
