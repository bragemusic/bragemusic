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

	token := "brg_v1_qrcbga8lH5cQsfBSyDbGn15A_zVEuvbw4oiGxbkyAvQ"

	_, err := client.NewFromToken(ctx, token, cCfg, slogHandler)
	if err != nil {
		logger.ErrorContext(ctx, "could not create client", "error", err.Error())
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh // blocks here until signal is received

	logger.Info("shutting down...")
}
