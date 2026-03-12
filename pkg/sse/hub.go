package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type EventHandler func(context.Context, types.SSEvent[any])

type ReqEvents struct {
	// AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	DeviceID uuid.UUID `path:"deviceID" description:"ID of your device"`
}

func (r ReqEvents) Validate() (validationMessages string, err error) {
	return "", nil
}

type eventEnvelope struct {
	event    types.SSEventBase
	deviceID *uuid.UUID
}

type client struct {
	deviceID uuid.UUID
	userID   uuid.UUID
	events   chan eventEnvelope
}

type Dispatcher interface {
	Broadcast(types.SSEventBase) error
	ActiveDevices(userID uuid.UUID) []uuid.UUID
	SendToDevice(deviceID uuid.UUID, ev types.SSEventBase) error
	SendToUser(userID uuid.UUID, ev types.SSEventBase) error
}

type Hub struct {
	db  database.DeviceFace
	log *slog.Logger

	clients map[*client]struct{}

	add       chan *client
	remove    chan *client
	broadcast chan eventEnvelope
}

func NewHub(db database.DeviceFace, slogHandler slog.Handler) *Hub {
	return &Hub{
		db:        db,
		log:       slog.New(slogHandler).With("service", "sse-hub"),
		clients:   make(map[*client]struct{}),
		add:       make(chan *client),
		remove:    make(chan *client),
		broadcast: make(chan eventEnvelope),
	}
}

func (h *Hub) Broadcast(ev types.SSEventBase) error {
	if len(h.clients) == 0 {
		return errors.New("no active clients")
	}
	h.broadcast <- eventEnvelope{
		event:    ev,
		deviceID: nil,
	}
	return nil
}

func (h *Hub) SendToDevice(deviceID uuid.UUID, ev types.SSEventBase) error {
	found := false
	for c := range h.clients {
		if c.deviceID == deviceID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("device not found")
	}

	h.broadcast <- eventEnvelope{
		event:    ev,
		deviceID: &deviceID,
	}
	return nil
}

func (h *Hub) SendToUser(userID uuid.UUID, ev types.SSEventBase) error {
	evEnvs := []eventEnvelope{}

	for c := range h.clients {
		if c.userID == userID {
			evEnvs = append(evEnvs, eventEnvelope{
				event:    ev,
				deviceID: &c.deviceID,
			})
		}
	}

	if len(evEnvs) == 0 {
		return errors.New("no user devices found")
	}

	for _, e := range evEnvs {
		h.broadcast <- e
	}

	return nil
}

func (h *Hub) ActiveDevices(userID uuid.UUID) (d []uuid.UUID) {
	for c := range h.clients {
		if c.userID == userID {
			d = append(d, c.deviceID)
		}
	}
	return d
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {

		case c := <-h.add:
			h.clients[c] = struct{}{}

		case c := <-h.remove:
			delete(h.clients, c)
			close(c.events)

		case e := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.events <- e:
				default:
					// client is too slow
					delete(h.clients, c)
					close(c.events)
				}
			}

		case <-ctx.Done():
			h.log.InfoContext(ctx, "recieved exit signal. Killing client connections")
			for c := range h.clients {
				delete(h.clients, c)
				close(c.events)
			}
			return
		}
	}
}

func (h *Hub) EventsHandler() routes.RouteFunc[ReqEvents, types.NoResponse] {
	return func(ctx context.Context, req ReqEvents, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		device, err := h.db.GetDeviceFromTokenAndDeviceID(ctx, req.DeviceID, tokenID)
		if err != nil {
			// FIXME: Bragerr device not registered
			return types.Response[types.NoResponse]{}, err
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			return types.Response[types.NoResponse]{}, errors.New("Streaming unsupported")
		}

		if err = h.db.UpdateDeviceLastSeen(ctx, "updated.ip.address.now", device.ID); err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		client := &client{
			deviceID: device.ID,
			userID:   user.ID,
			events:   make(chan eventEnvelope, 16),
		}

		h.log.InfoContext(ctx, "new client subscription", "device.name", device.Name, "user.email", user.Email)

		h.add <- client
		defer func() {
			select {
			case h.remove <- client:
			default:
			}
			h.log.InfoContext(ctx, "client disconnected", "device.name", device.Name, "user.email", user.Email)
		}()

		resp.Status = http.StatusOK

		for {
			select {

			case <-ctx.Done():
				return

			case e, ok := <-client.events:
				if !ok {
					// channel closed -> hub shutting down
					return
				}
				if e.deviceID == nil || *e.deviceID == device.ID {
					if err = sendEvent(w, flusher, e.event); err != nil {
						h.log.ErrorContext(ctx, "could not send event", "device.name", device.Name, "user.email", user.Email, "error", err.Error())
						return
					}
				}
			}
		}
	}
}

func sendEvent(w http.ResponseWriter, flusher http.Flusher, event types.SSEventBase) error {
	// Marshal the entire event to JSON
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Always use `data:` field; browsers will parse it as one JSON object
	fmt.Fprintf(w, "data: %s\n\n", jsonBytes)
	flusher.Flush()
	return nil
}
