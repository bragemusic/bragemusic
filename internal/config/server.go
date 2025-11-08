package config

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type ServerConfig struct {
	Port              string
	MusicDirPath      string
	ArtistArtsDirPath string
	AlbumArtsDirPath  string
	DBPath            string
	AcoustIDApiKey    string
	WikiEmail         string
}

func GetServerConfig() (sc ServerConfig, err error) {
	errs := []error{}

	sc.Port = envWithDefault("SERVER_PORT", "3000")

	sc.MusicDirPath, err = getEnv("MUSIC_DIR")
	errs = append(errs, err)

	imgDir, err := getEnv("IMG_DIR")
	errs = append(errs, err)

	sc.ArtistArtsDirPath = path.Join(imgDir, "artists")
	sc.AlbumArtsDirPath = path.Join(imgDir, "albums")

	serverDir, err := getEnv("SERVER_DIR")
	errs = append(errs, err)

	sc.DBPath = path.Join(serverDir, "data.db")

	sc.AcoustIDApiKey, err = getEnv("ACOUST_ID_APIKEY")
	errs = append(errs, err)

	sc.WikiEmail, err = getEnv("WIKI_EMAIL")
	errs = append(errs, err)

	return sc, errors.Join(errs...)
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("env %s not set", key)
	}
	return val, nil
}

func envWithDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
