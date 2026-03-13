package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/bragemusic/core/pkg/client"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
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
		PlayerName:      "Lucas Daemon",
		ClientType:      types.DeviceTypeStreaming,
		ClientInterface: types.DeviceInterfaceDaemon,
		StateFilePath:   utils.Ptr("/home/lucas/dev/brage/daemondata"),
	}
	c, err := client.New(ctx, cCfg, slogHandler)
	if err != nil {
		logger.ErrorContext(ctx, "could not create client", "error", err.Error())
	}

	token := "brg_v1_yEH8oLaY7CET0fsdpBExPkaBcWKpgIWa_G9IWH573AI"

	user, err := c.LoginToken(ctx, token)
	if err != nil {
		logger.ErrorContext(ctx, "could not log in", "error", err.Error())
	}

	logger.InfoContext(ctx, "successfully logged in", "user.name", user.Username, "user.email", user.Email)

	err = c.StartPlayerWithAlbum(ctx, uuid.Must(uuid.FromString("773e8a2a-9627-4eb0-8929-1004498375cb")), 0)
	if err != nil {
		logger.ErrorContext(ctx, "could not play", "error", err.Error())
	}

	for {
	}
}
