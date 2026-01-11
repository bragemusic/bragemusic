package imagemagick

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type ImageSize int

const (
	Size320  ImageSize = 320
	Size640  ImageSize = 640
	Size1024 ImageSize = 1024
	Size1600 ImageSize = 1600
	Size2400 ImageSize = 2400
)

type ImageMagick struct {
	log *slog.Logger
}

func (i ImageMagick) Resize(ctx context.Context, inputFile, outputFile string, size ImageSize) error {
	sizeStr := fmt.Sprintf("%dx%d^", size, size)

	cmd := exec.CommandContext(ctx, "magick", inputFile, "-resize", sizeStr, "-strip", "-interlace", "Plane", "-quality", "82", outputFile)
	if output, err := cmd.Output(); err != nil {
		i.log.ErrorContext(ctx, "could not run magick command", "error", err.Error(), "output", output, "args", cmd.Args)
		return err
	}

	return nil
}

func (i ImageMagick) ResizeAll(ctx context.Context, inputFile, outputFolder string) error {
	if err := os.MkdirAll(outputFolder, 0o755); err != nil {
		i.log.ErrorContext(ctx, "could not create image dir", "error", err.Error())
		return err
	}

	for _, size := range []ImageSize{Size320, Size640, Size1024, Size1600, Size2400} {
		outFilename := filepath.Join(outputFolder, fmt.Sprintf("%d.jpg", size))
		if err := i.Resize(ctx, inputFile, outFilename, size); err != nil {
			return err
		}
	}

	return nil
}

func New(slogHandler slog.Handler) (ImageMagick, error) {
	_, err := exec.LookPath("magick")
	if err != nil {
		return ImageMagick{}, errors.New("'magick' command not found on the computer")
	}

	return ImageMagick{
		log: slog.New(slogHandler).With("service", "imagemagick"),
	}, nil
}
