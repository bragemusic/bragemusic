package utils

import (
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	discTrackPatterns = []*regexp.Regexp{
		// Disc 1 - 03 Track, Disc1_03 Track
		regexp.MustCompile(`(?i)disc\s*(\d{1,2})\D+(\d{1,3})`),

		// CD2-07 Track, CD 2 - 07
		regexp.MustCompile(`(?i)cd\s*(\d{1,2})\D+(\d{1,3})`),

		// 1-03 Track, 1.03 Track, 1_03 Track
		regexp.MustCompile(`(?i)^\s*(\d{1,2})[._-](\d{1,3})\b`),

		// [2] [07] Track, (2)(07)
		regexp.MustCompile(`(?i)[\[\(]\s*(\d{1,2})\s*[\]\)]\s*[\[\(]\s*(\d{1,3})\s*[\]\)]`),

		// Disc only missing but two numbers at start: 02 07 Track
		regexp.MustCompile(`(?i)^\s*(\d{1,2})\s+(\d{1,3})\b`),
	}

	trackOnlyPatterns = []*regexp.Regexp{
		// 03 Track, 03-Track, 03.Track
		regexp.MustCompile(`(?i)^\s*(\d{1,3})[\s._-]+`),

		// [03] Track
		regexp.MustCompile(`(?i)^[\[\(]\s*(\d{1,3})\s*[\]\)]`),

		// Track 03
		regexp.MustCompile(`(?i)\b(\d{1,3})\s*$`),
	}
)

// ExtractDiscAndTrack returns (disc, track, ok)
func ExtractDiscAndTrack(filename string) (int, int, bool) {
	name := filepath.Base(filename)
	ext := filepath.Ext(name)
	name = name[:len(name)-len(ext)]

	// Try disc + track patterns first
	for _, re := range discTrackPatterns {
		if m := re.FindStringSubmatch(name); len(m) == 3 {
			disc, err1 := strconv.Atoi(m[1])
			track, err2 := strconv.Atoi(m[2])
			if err1 == nil && err2 == nil && validDiscTrack(disc, track) {
				return disc, track, true
			}
		}
	}

	// Fallback: track only (disc = 1)
	for _, re := range trackOnlyPatterns {
		if m := re.FindStringSubmatch(name); len(m) == 2 {
			track, err := strconv.Atoi(m[1])
			if err == nil && track > 0 && track < 1000 {
				return 1, track, true
			}
		}
	}

	return 0, 0, false
}

func validDiscTrack(disc, track int) bool {
	return disc > 0 && disc < 100 && track > 0 && track < 1000
}
