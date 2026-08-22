package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"github.com/bragemusic/bragemusic/pkg/client"
	"github.com/bragemusic/bragemusic/pkg/types"
)

type General struct {
	PlayerName         string           `toml:"player_name" desc:"The name of your device. This is what will be shown in the connect funtionality."`
	PlayerIcon         types.DeviceIcon `toml:"player_icon" desc:"The icon of your player will have. Choose from 'laptop', 'computer', 'phone', 'speaker', 'tv', 'generic'. Defaults to 'laptop'."`
	DisableTransitions bool             `toml:"disable_transitions" desc:"Set this to true if you see weird artefacts on popups, overlays and more. Will disable all transitions."`
	ClientType         types.DeviceType `toml:"client_type" desc:"Either sync or streaming. Sync client is syncing all files to the client and can run offline. The streaming client streams everything from the server."`
	LogFormat          LogFormat        `toml:"log_format" desc:"The format of the logs printed to the stderr. Defaults to 'pretty'"`
	LogLevel           string           `toml:"log_level" desc:"Set the log level of the client. Defaults to 'INFO'."`
}

func (g General) Verify() (errs []VerificationError) {
	switch g.PlayerIcon {
	case types.DeviceIconLaptop, types.DeviceIconComputer, types.DeviceIconPhone, types.DeviceIconTV, types.DeviceIconSpeaker, types.DeviceIconGeneric:
	default:
		errs = append(errs, VerificationError{
			Parameter: "General.PlayerIcon",
			Error:     fmt.Sprintf("unknown player icon name '%s'. Choose from Choose from 'laptop', 'computer', 'phone', 'speaker', 'tv', 'generic'", g.PlayerIcon),
		})
	}

	switch g.ClientType {
	case types.DeviceTypeStreaming, types.DeviceTypeSync:
	default:
		errs = append(errs, VerificationError{
			Parameter: "General.ClientType",
			Error:     fmt.Sprintf("unknown client type '%s'. Choose from Choose from 'sync', 'streaming'", g.ClientType),
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

type ClientPaths struct {
	ConfigDir string `toml:"config_dir"`
	ImageDir  string `toml:"image_dir"`
	MusicDir  string `toml:"music_dir"`
}

type Server struct {
	BaseUrl string `toml:"base_url"`
}

func (s Server) Verify() (errs []VerificationError) {
	_, err := url.ParseRequestURI(s.BaseUrl)
	if err != nil {
		errs = append(errs, VerificationError{
			Parameter: "Server.BaseUrl",
			Error:     fmt.Sprintf("'%s' is not a valid server url", s.BaseUrl),
		})
	}

	return errs
}

type Auth struct {
	Token string `toml:"token" desc:"Token to authenticate with the server. Required if running in daemon mode."`
}

type ClientConfig struct {
	Auth    Auth                   `toml:"auth"`
	General General                `toml:"general"`
	Paths   ClientPaths            `toml:"paths"`
	Server  Server                 `toml:"server"`
	Themes  map[string]types.Theme `toml:"theme"`
}

func (c ClientConfig) Verify() (errs []VerificationError) {
	errs = append(errs, c.General.Verify()...)
	errs = append(errs, c.Server.Verify()...)
	// TODO: Add verification
	return errs
}

func (c ClientConfig) ClientConfig() client.Config {
	cCfg := client.Config{
		ServerBaseURL: c.Server.BaseUrl,
		MusicDirPath:  c.Paths.MusicDir,
		ConfigPath:    c.Paths.ConfigDir,
		ImagePath:     c.Paths.ImageDir,
		PlayerName:    c.General.PlayerName,
		ClientType:    c.General.ClientType,
		ClientIcon:    c.General.PlayerIcon,
		// FIXME: Dont hardcode
		ClientInterface: types.DeviceInterfaceDesktop,
	}
	return cCfg
}

var defaultClientConfig = ClientConfig{
	General: General{
		PlayerName: "brage-music",
		PlayerIcon: types.DeviceIconLaptop,
		ClientType: types.DeviceTypeStreaming,
		LogLevel:   "INFO",
		LogFormat:  LogFormatPretty,
	},
	Paths: ClientPaths{
		ConfigDir: filepath.Join(xdg.DataHome, "brage", "config"),
		ImageDir:  filepath.Join(xdg.DataHome, "brage", "img"),
		MusicDir:  filepath.Join(xdg.DataHome, "brage", "music"),
	},
	Server: Server{},
}

func readCfgFile(filename string, cfg ClientConfig) (ClientConfig, error) {
	if _, err := os.Stat(filename); err != nil {
		return cfg, err
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	_, err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// func verify2(cfg Config) (errs []VerificationError) {
// 	errs = append(errs, cfg.General.Verify()...)
// 	errs = append(errs, cfg.Server.Verify()...)
// 	// TODO: Add verification
// 	return errs
// }

func MakeSureUserHasConfig() error {
	cfgFilepath := filepath.Join(xdg.ConfigHome, "brage", "config.toml")

	_, err := os.Stat(cfgFilepath)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(cfgFilepath), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(cfgFilepath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = toml.NewEncoder(f).Encode(defaultClientConfig); err != nil {
		return err
	}

	return nil
}

func GetClientConfig(logger *slog.Logger) (ClientConfig, error) {
	c := defaultClientConfig

	cfgFilepath := filepath.Join(xdg.ConfigHome, "brage", "config.toml")

	c, err := readCfgFile(cfgFilepath, c)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return ClientConfig{}, err
		}
	}

	err = addFromEnv(&c)
	if err != nil {
		return ClientConfig{}, err
	}

	if err = verify(c, logger); err != nil {
		return ClientConfig{}, err
	}

	return c, nil
}

func ClientMdDocs() (string, error) {
	docs, err := generateConfigDocs(&ClientConfig{Themes: map[string]types.Theme{"theme_name": {}}})
	if err != nil {
		return "", err
	}

	header := `
## Client Config
The following configuration is used in the local client. The daemon and the GUI are both accepting the same config.
It is stored in ~/.config/brage/config.toml or use it with ENV variables.
`
	return header + docs, nil
}
