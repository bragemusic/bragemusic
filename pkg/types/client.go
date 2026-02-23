package types

import "log/slog"

type ClientConfig struct {
	ConfigPath    string
	ImagePath     string
	MusicDirPath  string
	PlayerName    string
	ServerBaseURL string
}

type ClientMessage struct {
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Level   slog.Level `json:"level"`
}
