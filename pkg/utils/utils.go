package utils

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/bragemusic/bragemusic/pkg/acoustid"
	"github.com/dhowden/tag"
)

var mbIDRe = regexp.MustCompile(`mbid#([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)

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

func CmpReleaseDates(d1, d2 acoustid.Date) int {
	d1 = fixDate(d1)
	d2 = fixDate(d2)

	// Compare year
	if d1.Year < d2.Year {
		return -1
	}

	if d1.Year > d2.Year {
		return 1
	}

	// Compare month
	if d1.Month < d2.Month {
		return -1
	}

	if d1.Month > d2.Month {
		return 1
	}

	// Compare day
	if d1.Day < d2.Day {
		return -1
	}

	if d1.Day > d2.Day {
		return 1
	}

	// Equal
	return 0
}

func fixDate(d acoustid.Date) acoustid.Date {
	if d.Year == 0 {
		d.Year = 9999
	}

	if d.Month == 0 {
		d.Month = 12
	}

	if d.Day == 0 {
		d.Day = 31
	}

	return d
}

func FileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	checksum := hasher.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func ParseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	return level, err
}

func MatchMbID(s string) (found bool, mbid string) {
	matches := mbIDRe.FindStringSubmatch(s)
	if matches == nil {
		return false, ""
	}

	mbid = matches[1]

	return true, mbid
}
