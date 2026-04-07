package types

import (
	"log/slog"
)

type (
	ClientEvent string
)

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

const (
	ClientEventMsgErr     ClientEvent = "msg.error"
	ClientEventMsgSuccess ClientEvent = "msg.success"
	ClientEventMsgInfo    ClientEvent = "msg.info"
	ClientEventMsgWarn    ClientEvent = "msg.warn"

	ClientEventServerOnline ClientEvent = "server.status"

	ClientEventSyncInProgress ClientEvent = "sync.inprogress"

	ClientEventEntitiesUpdated ClientEvent = "entities.updated"

	ClientEventPlayerContextChange  ClientEvent = "player.contextchange"
	ClientEventPlayerPlaybackChange ClientEvent = "player.playbackchange"

	ClientEventUserUpdated ClientEvent = "user.updated"

	ClientEventAuthLoggedIn ClientEvent = "auth.loggedin"
)
