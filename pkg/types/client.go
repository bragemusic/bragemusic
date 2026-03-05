package types

import (
	"log/slog"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	ClientEvent string
	ClientType  string
)

const (
	ClientTypeStreaming ClientType = "streaming"
	ClientTypeSync      ClientType = "sync"
	ClientTypeDaemon    ClientType = "daemon"
	ClientTypeWeb       ClientType = "web"
)

// TODO:UNIQUE(token_id)
type Client struct {
	ID               uuid.UUID
	Name             string
	Type             ClientType
	TokenID          uuid.UUID
	SupportsPlayback bool
	Platform         string
	Version          string
	LastIP           string
	LastSeen         time.Time
}

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

	ClientEventPlayerContextChange ClientEvent = "player.contextchange"

	ClientEventUserUpdated ClientEvent = "user.updated"
)
