package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bragemusic/core/pkg/client"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/lmittmann/tint"
)

func main() {
	ctx := context.Background()

	slogHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	logger := slog.New(slogHandler)

	cCfg := client.Config{
		ServerBaseURL:   "http://localhost:3000",
		MusicDirPath:    "/home/lucas/dev/brage/client-data/music",
		ConfigPath:      "/home/lucas/dev/brage/client-data/config",
		ImagePath:       "/home/lucas/dev/brage/client-data/img",
		PlayerName:      "Stereo",
		ClientIcon:      types.DeviceIconSpeaker,
		ClientType:      types.DeviceTypeStreaming,
		ClientInterface: types.DeviceInterfaceDaemon,
		StateFilePath:   utils.Ptr("/home/lucas/dev/brage/daemondata"),
	}
	c, err := client.New(ctx, cCfg, slogHandler)
	if err != nil {
		logger.ErrorContext(ctx, "could not create client", "error", err.Error())
	}

	token := "brg_v1_VtfPBMEtN4otYzPslFFtxr2CcYcZSYveaK5v5txHI9Q"

	user, err := c.LoginToken(ctx, token)
	if err != nil {
		logger.ErrorContext(ctx, "could not log in", "error", err.Error())
		return
	}

	logger.InfoContext(ctx, "successfully logged in", "user.name", user.Username, "user.email", user.Email)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh // blocks here until signal is received

	logger.Info("shutting down...")
}
