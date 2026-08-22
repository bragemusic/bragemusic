package migrations

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
	dbMigrations "github.com/bragemusic/bragemusic/db"
)

func Migrate(ctx context.Context, databasePath string, slogHandler slog.Handler) error {
	logger := slog.New(slogHandler).With("service", "migrations")

	sw := slogWriter{
		logger: logger,
		ctx:    ctx,
		level:  slog.LevelInfo,
	}

	u, err := url.Parse(fmt.Sprintf("sqlite:%s", databasePath))
	if err != nil {
		return err
	}

	db := dbmate.New(u)
	db.FS = dbMigrations.FS
	db.MigrationsDir = []string{"migrations"}
	db.AutoDumpSchema = false
	db.Log = sw

	err = db.CreateAndMigrate()
	if err != nil {
		return err
	}

	return nil
}

type slogWriter struct {
	logger *slog.Logger
	level  slog.Level
	ctx    context.Context
}

func (w slogWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	msg = strings.TrimRight(msg, "\n")

	w.logger.Log(w.ctx, w.level, msg)
	return len(p), nil
}
