package filetx

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var logger *slog.Logger

type FileTx struct {
	ctx       context.Context
	rollbacks []func() error
}

func (ftx *FileTx) CopyFile(from, to string) error {
	logger.DebugContext(ftx.ctx, "copying file", "src", from, "dst", to)

	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		return err
	}

	src, err := os.OpenFile(from, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	ftx.rollbacks = append(ftx.rollbacks, func() error {
		logger.DebugContext(ftx.ctx, "removing file", "filename", to)
		return os.Remove(to)
	})

	return nil
}

func (ftx *FileTx) Commit() error {
	logger.DebugContext(ftx.ctx, "commiting file transaction")
	ftx.rollbacks = nil
	return nil
}

func (ftx *FileTx) Rollback() error {
	if len(ftx.rollbacks) == 0 {
		return nil
	}

	logger.DebugContext(ftx.ctx, "rolling back file transaction")
	errs := []error{}
	for _, f := range ftx.rollbacks {
		errs = append(errs, f())
	}

	return errors.Join(errs...)
}

func Init(slogHandler slog.Handler) {
	logger = slog.New(slogHandler).With("service", "filetx")
}

func Begin(ctx context.Context) (FileTx, error) {
	if logger == nil {
		return FileTx{}, errors.New("filetx must be initialized before use")
	}
	return FileTx{
		ctx: ctx,
	}, nil
}
