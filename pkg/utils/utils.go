package utils

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dhowden/tag"
)

func Ptr[T any](t T) *T {
	return &t
}

func SaveID3Image(ctx context.Context, img tag.Picture, filename string) error {
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

func HighestCount[T comparable](ss []T) T {
	uss := slices.Compact(ss)
	counts := map[T]int{}

	for _, s := range uss {
		counts[s] = 0
		for _, a := range ss {
			if a == s {
				counts[s]++
			}
		}
	}

	bestName := ss[0]
	bestCnt := counts[ss[0]]
	for an, ac := range counts {
		if ac > bestCnt {
			bestName = an
			bestCnt = ac
		}
	}

	return bestName
}
