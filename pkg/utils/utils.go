package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

func Ptr[T any](t T) *T {
	return &t
}

func SaveID3Image(ctx context.Context, img *tag.Picture, filename string) error {
	if img == nil {
		return errors.New("no image data in ID3 tag")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(img.Data)

	return nil
}

func GenerateAlbumFolderPath(artist, album string) string {
	artist = strings.ReplaceAll(artist, " ", "_")
	album = strings.ReplaceAll(album, " ", "_")

	return path.Join(artist, album)
}

func GenerateTrackPath(discNumber, trackNumber int, trackTitle string, format tag.FileType, albumFolder string) (string, error) {
	if format == "" {
		return "", fmt.Errorf("unknown fileformat for track '%s'", trackTitle)
	}

	trackTitle = strings.ReplaceAll(trackTitle, " ", "_")
	trackTitle = strings.ReplaceAll(trackTitle, "/", "_")

	filename := fmt.Sprintf("%02d-%02d-%s.%s", discNumber, trackNumber, trackTitle, strings.ToLower(string(format)))

	return filepath.Join(albumFolder, filename), nil
}
