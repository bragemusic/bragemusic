package gui

import (
	"context"
	"embed"
	"log/slog"

	"github.com/bragemusic/core/pkg/client"
	"github.com/bragemusic/core/pkg/config"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/bragemusic/core/internal/app"
	"github.com/bragemusic/core/internal/assethandler"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

func Run(ctx context.Context, logger *slog.Logger, cfg config.ClientConfig, assets embed.FS) {
	logger.Info("Starting Brage Music GUI")

	// FIXME: Sync Daemon needs to use wails app context to be able to finish a sync if user quits
	aC := client.NewAuthClient(ctx, cfg.ClientConfig(), logger.Handler())

	app := app.New(aC, cfg, logger.Handler())

	lo := linux.Options{
		WebviewGpuPolicy: linux.WebviewGpuPolicyAlways,
	}
	// Create application with options
	err := wails.Run(&options.App{
		Title: "Brage Music",
		// Width:  1024,
		// Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: utils.Ptr(assethandler.New(cfg, aC.ServerClient())),
		},
		// BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:  app.Startup,
		OnShutdown: app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Linux: &lo,
	})
	if err != nil {
		logger.Error("could not run main app", "error", err.Error())
		return
	}
}
