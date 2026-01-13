package mediamanager

import (
	"log/slog"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/imagemagick"
)

type MediaManager struct {
	log      *slog.Logger
	db       database.DatabaseFace
	im       *imagemagick.ImageMagick
	musicDir string
	imageDir string
}

func New(slogHandler slog.Handler, db database.DatabaseFace, im *imagemagick.ImageMagick, musicDir, imageDir string) MediaManager {
	return MediaManager{
		log:      slog.New(slogHandler).With("service", "media-manager"),
		db:       db,
		im:       im,
		musicDir: musicDir,
		imageDir: imageDir,
	}
}
