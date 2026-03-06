package sse

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type ReqSubscribe struct {
	// AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	ID               uuid.UUID             `json:"id"`
	Name             string                `json:"name"`
	Type             types.DeviceType      `json:"type"`
	Interface        types.DeviceInterface `json:"interface"`
	SupportsPlayback bool                  `json:"supports_playback"`
	Platform         string                `json:"platform"`
	Version          string                `json:"version"`
}

func (r ReqSubscribe) Validate() (validationMessages string, err error) {
	return "", nil
}

type Event struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	Data any       `json:"data"`
}

type client struct {
	id     uuid.UUID
	events chan Event
}

type Hub struct {
	db  database.DeviceFace
	log *slog.Logger

	clients map[*client]struct{}

	add       chan *client
	remove    chan *client
	broadcast chan Event
}

func NewHub(db database.DeviceFace, slogHandler slog.Handler) *Hub {
	return &Hub{
		db:        db,
		log:       slog.New(slogHandler).With("service", "sse-hub"),
		clients:   make(map[*client]struct{}),
		add:       make(chan *client),
		remove:    make(chan *client),
		broadcast: make(chan Event),
	}
}

func (h *Hub) Broadcast(ev Event) error {
	if len(h.clients) == 0 {
		return errors.New("no active clients")
	}
	h.broadcast <- ev
	return nil
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

func (h *Hub) Handler() routes.RouteFunc[ReqSubscribe, types.NoResponse] {
	return func(ctx context.Context, req ReqSubscribe, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			return types.Response[types.NoResponse]{}, err
		}

		device, err := h.db.GetDeviceFromTokenID(ctx, tokenID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return types.Response[types.NoResponse]{}, err
			}
			device = types.Device{
				ID:               req.ID,
				Name:             req.Name,
				Type:             req.Type,
				Interface:        req.Interface,
				TokenID:          tokenID,
				SupportsPlayback: req.SupportsPlayback,
				Platform:         req.Platform,
				Version:          req.Version,
				LastIP:           "will.fix.this.later",
				LastSeen:         time.Now(),
			}
			if err = h.db.AddDevice(ctx, device); err != nil {
				return types.Response[types.NoResponse]{}, err
			}
		} else {
			device.Name = req.Name
			device.Type = req.Type
			device.Interface = req.Interface
			device.SupportsPlayback = req.SupportsPlayback
			device.Platform = req.Platform
			device.Version = req.Version
			device.LastIP = "will.fix.this.later"
			device.LastSeen = time.Now()

			if err = h.db.UpdateDevice(ctx, device); err != nil {
				return types.Response[types.NoResponse]{}, err
			}
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			return types.Response[types.NoResponse]{}, errors.New("Streaming unsupported")
		}

		client := &client{
			id:     device.ID,
			events: make(chan Event, 16),
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
				if err = sendEvent(w, flusher, e); err != nil {
					h.log.ErrorContext(ctx, "could not send event", "device.name", device.Name, "user.email", user.Email, "error", err.Error())
					return
				}
			}
		}
	}
}

func sendEvent(w http.ResponseWriter, flusher http.Flusher, event Event) error {
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
