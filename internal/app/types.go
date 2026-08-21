package app

type (
	Event string
)

const (
	EventMsgErr     Event = "msg.error"
	EventMsgSuccess Event = "msg.success"
	EventMsgInfo    Event = "msg.info"
	EventMsgWarn    Event = "msg.warn"

	EventServerOnline   Event = "server.status"
	EventSyncInProgress Event = "sync.inprogress"

	EventEntitiesUpdated Event = "entities.updated"

	EventPlayerContextChange Event = "player.contextchange"

	EventUserUpdated Event = "user.updated"
)

type Message struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
