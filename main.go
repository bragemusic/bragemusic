package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	chigo "github.com/UltiRequiem/chigo/pkg"
	"github.com/alecthomas/kong"
	"github.com/bragemusic/core/internal/daemon"
	"github.com/bragemusic/core/internal/gui"
	"github.com/bragemusic/core/pkg/config"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/common-nighthawk/go-figure"
	"github.com/lmittmann/tint"
)

//go:embed all:wails/frontend/dist
var assets embed.FS

var cmdDesc = chigo.Colorize(figure.NewFigure("Brage Music", "speed", true).String()) + `
A cross-platform desktop client for BrageMusic with a built-in graphical interface
and optional daemon mode. Connect to and manage your BrageMusic server, browse your library,
create playlists and smart playlists, and control playback. Run it in the background for
lightweight system integration, or launch the full GUI for a complete desktop music experience.
`

type CLI struct {
	GUI        struct{} `cmd:"" aliases:"g," default:"1" help:"Start the application GUI. (Default)."`
	Daemon     struct{} `cmd:"" aliases:"d," help:"Start the application in daemon mode. No GUI, only terminal."`
	ConfigDocs struct{} `cmd:"" help:"View configuration documentation."`
}

func main() {
	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	logger := slog.New(slogHandler)

	err := config.MakeSureUserHasConfig()
	if err != nil {
		logger.Error("could not generate config", "error", err.Error())
		return
	}

	cfg, err := config.GetClientConfig(logger)
	if err != nil {
		logger.Error("could not get config", "error", err.Error())
		return
	}

	logLevel, err := utils.ParseLogLevel(cfg.General.LogLevel)
	if err != nil {
		logger.Error("could not parse log level. Falling back to DEBUG", "error", err.Error())
	} else {
		switch cfg.General.LogFormat {
		case config.LogFormatPretty:
			slogHandler = tint.NewHandler(os.Stderr, &tint.Options{
				Level:      logLevel,
				TimeFormat: time.TimeOnly,
			})
		case config.LogFormatJson:
			slogHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: logLevel,
			})
		default:
			logger.Error("unknown log format", "log_format", cfg.General.LogFormat)
		}

		logger = slog.New(slogHandler)
	}

	kctx := kong.Parse(&CLI{},
		kong.Name("bragemusic"),
		kong.Description(cmdDesc),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppDescFormat: true,
			Compact:         false,
			Summary:         false,
		}),
	)

	ctx := context.Background()

	switch kctx.Command() {
	case "gui":
		gui.Run(ctx, logger, cfg, assets)
	case "daemon":
		daemon.Run(ctx, logger, cfg)
	case "config-docs":
		source, err := config.ClientMdDocs()
		if err != nil {
			logger.ErrorContext(ctx, "could not render config docs", "error", err.Error())
			return
		}
		result := markdown.Render(string(source), 80, 6)
		fmt.Println(string(result))
		return
	default:
		logger.ErrorContext(ctx, "command not found", "command", kctx.Command())
		return
	}
}
