package daemon

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bragemusic/core/pkg/client"
	"github.com/bragemusic/core/pkg/config"
	"github.com/bragemusic/core/pkg/types"
)

func Run(ctx context.Context, logger *slog.Logger, cfg config.ClientConfig) {
	token := cfg.Auth.Token
	if token == "" {
		logger.ErrorContext(ctx, "no token provided. Cannot start application")
		return
	}

	cCfg := cfg.ClientConfig()
	cCfg.ClientType = types.DeviceTypeStreaming
	cCfg.ClientInterface = types.DeviceInterfaceDaemon

	_, err := client.NewFromToken(ctx, token, cCfg, logger.Handler())
	if err != nil {
		logger.ErrorContext(ctx, "could not create client", "error", err.Error())
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh // blocks here until signal is received

	logger.Info("shutting down...")
}
