package config

import (
	"fmt"
	"log/slog"

	"github.com/adhocore/gronx"
)

type Paths struct {
	ConfigDir       string `toml:"config_dir" desc:"Dir where server files are stored."`
	ImageDir        string `toml:"image_dir" desc:"Dir where image assets are stored."`
	MusicDir        string `toml:"music_dir" desc:"Dir where music files are stored."`
	ImportDir       string `toml:"import_dir" desc:"Dir where imported albums and tracks will be saved before processing."`
	ManualImportDir string `toml:"manual_import_dir" desc:"Dir where the importer is looking for manually added bulk imports."`
	BackupImportDir string `toml:"backup_import_dir" desc:"Dir where imported albums and tracks will be saved after processing."`
}

func (p Paths) Verify() (errs []VerificationError) {
	if p.ImageDir == "" {
		errs = append(errs, VerificationError{
			Parameter: "Paths.ImageDir",
			Error:     "parameter not set",
		})
	}

	if p.MusicDir == "" {
		errs = append(errs, VerificationError{
			Parameter: "Paths.MusicDir",
			Error:     "parameter not set",
		})
	}

	if p.ConfigDir == "" {
		errs = append(errs, VerificationError{
			Parameter: "Paths.ConfigDir",
			Error:     "parameter not set",
		})
	}

	if p.ImportDir == "" {
		errs = append(errs, VerificationError{
			Parameter: "Paths.ImportDir",
			Error:     "parameter not set",
		})
	}

	if p.ManualImportDir == "" {
		errs = append(errs, VerificationError{
			Parameter: "Paths.ManualImportDir",
			Error:     "parameter not set",
		})
	}

	if p.BackupImportDir == "" {
		errs = append(errs, VerificationError{
			Parameter: "Paths.BackupImportDir",
			Error:     "parameter not set",
		})
	}

	return errs
}

type Admin struct {
	Email    string `toml:"email" desc:"Default user, with admin rights, email. Defaults to 'admin@example.com'"`
	Username string `toml:"username" desc:"Default user, with admin rights, username. Defaults to 'admin'"`
	Password string `toml:"password" desc:"Default user, with admin rights, password. Defaults to 'password'"`
}

type AcoustID struct {
	ApiKey string `toml:"api_key" desc:"API key to Acoust ID. Used to identify the files to not rely solely on ID3."`
}

func (a AcoustID) Verify() (errs []VerificationError) {
	if a.ApiKey == "" {
		errs = append(errs, VerificationError{
			Parameter: "AcoustID.ApiKey",
			Error:     "must be set to analyse tracks",
		})
	}

	return errs
}

type Analyser struct {
	BaseURL string `toml:"base_url" desc:"URL to the analysis service. If left blank, no analysis will be performed."`
}

type Wikipedia struct {
	Email string `toml:"email" desc:"Used against the wikipedia API. They require a valid email to make sure you behave."`
}

func (w Wikipedia) Verify() (errs []VerificationError) {
	if w.Email == "" {
		errs = append(errs, VerificationError{
			Parameter: "Wikipedia.Email",
			Error:     "must be set to be able to access wikipedia",
		})
	}

	return errs
}

type Jobs struct {
	Importer      string `toml:"importer" desc:"How often the importer will look for new media files. Cron expression. Defaults to '*/3 * * * *'"`
	MetaSyncer    string `toml:"meta_syncer" desc:"How often the meta-syncer will sync the needed metadata. Cron expression. Defaults to '*/3 * * * *'"`
	SearchItems   string `toml:"search_items" desc:"How often the search items will be updated. Cron expression. Defaults to '*/3 * * * *'"`
	TrackAnalysis string `toml:"analyser" desc:"How often the track analysis items will be updated. Cron expression. Defaults to '*/3 * * * *'"`
	TokenCleanup  string `toml:"token_cleanup" desc:"How often token cleanup will be performed. Cron expression. Defaults to '*/10 * * * *'"`
}

func (j Jobs) Verify() (errs []VerificationError) {
	if !gronx.IsValid(j.Importer) {
		errs = append(errs, VerificationError{
			Parameter: "Jobs.Importer",
			Error:     "not a valid Cron Expression",
		})
	}

	if !gronx.IsValid(j.MetaSyncer) {
		errs = append(errs, VerificationError{
			Parameter: "Jobs.MetaSyncer",
			Error:     "not a valid Cron Expression",
		})
	}

	if !gronx.IsValid(j.SearchItems) {
		errs = append(errs, VerificationError{
			Parameter: "Jobs.SearchItems",
			Error:     "not a valid Cron Expression",
		})
	}

	if !gronx.IsValid(j.TrackAnalysis) {
		errs = append(errs, VerificationError{
			Parameter: "Jobs.TrackAnalysis",
			Error:     "not a valid Cron Expression",
		})
	}

	if !gronx.IsValid(j.TokenCleanup) {
		errs = append(errs, VerificationError{
			Parameter: "Jobs.TokenCleanup",
			Error:     "not a valid Cron Expression",
		})
	}

	return errs
}

type ServerGeneral struct {
	Name      string    `toml:"name" desc:"Name of the server. Defaults to 'Brage Music Server'"`
	Port      int       `toml:"port" desc:"Port of the server. Defaults to 3000."`
	LogFormat LogFormat `toml:"log_format" desc:"The format of the logs printed to the stderr. Defaults to 'json'"`
	LogLevel  string    `toml:"log_level" desc:"Set the log level of the client. Defaults to 'INFO'."`
}

func (g ServerGeneral) Verify() (errs []VerificationError) {
	if g.Port < 1024 || g.Port > 49151 {
		errs = append(errs, VerificationError{
			Parameter: "General.Port",
			Error:     fmt.Sprintf("port %d not allowed. Use in range 1024-49151", g.Port),
		})
	}

	switch g.LogFormat {
	case LogFormatJson, LogFormatPretty:
	default:
		errs = append(errs, VerificationError{
			Parameter: "General.LogFormat",
			Error:     fmt.Sprintf("unknown log format '%s'. Choose from 'pretty', 'json'", g.LogFormat),
		})
	}
	return errs
}

type ServerConfig struct {
	AcoustID  AcoustID      `toml:"acoust_id"`
	Admin     Admin         `toml:"admin"`
	Analyser  Analyser      `toml:"analyser"`
	General   ServerGeneral `toml:"general"`
	Jobs      Jobs          `toml:"jobs"`
	Paths     Paths         `toml:"paths"`
	Wikipedia Wikipedia     `toml:"wikipedia"`
}

func (s ServerConfig) Verify() (errs []VerificationError) {
	errs = append(errs, s.AcoustID.Verify()...)
	errs = append(errs, s.General.Verify()...)
	errs = append(errs, s.Jobs.Verify()...)
	errs = append(errs, s.Paths.Verify()...)
	errs = append(errs, s.Wikipedia.Verify()...)

	// TODO: Add more verification
	return errs
}

var defaultServerConfig = ServerConfig{
	Admin: Admin{
		Email:    "admin@example.com",
		Username: "admin",
		Password: "password",
	},
	Paths: Paths{},
	Jobs: Jobs{
		Importer:      "*/3 * * * *",
		MetaSyncer:    "*/3 * * * *",
		SearchItems:   "*/3 * * * *",
		TrackAnalysis: "*/3 * * * *",
		TokenCleanup:  "*/10 * * * *",
	},
	General: ServerGeneral{
		Name:      "Brage Music Server",
		Port:      3000,
		LogFormat: LogFormatJson,
		LogLevel:  "INFO",
	},
}

func GetServerConfig(logger *slog.Logger) (ServerConfig, error) {
	c := defaultServerConfig

	err := addFromEnv(&c)
	if err != nil {
		return ServerConfig{}, err
	}

	if err = verify(c, logger); err != nil {
		return ServerConfig{}, err
	}

	return c, nil
}

func ServerMdDocs() (string, error) {
	docs, err := generateConfigDocs(&ServerConfig{})
	if err != nil {
		return "", err
	}

	header := `
## Server Config
This config is used to define the server. It only accepts ENV values at the moment.
`
	return header + docs, nil
}
