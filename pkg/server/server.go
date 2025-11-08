package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	log *slog.Logger

	db database.DatabaseFace
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", s.healthz())

	r.Get("/artists", s.listArtists())
	r.Get("/artists/{artistID}/albums", s.listAlbums())

	return r
}

func (s Server) healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := s.writeJSON(ctx, w, Healthz{
			Application: "brage-server",
			Version:     "v0.0.1",
			Status:      HealthzRunning,
		})
		if err != nil {
			s.handleErr(ctx, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s Server) writeJSON(ctx context.Context, w http.ResponseWriter, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)
	return err
}

func (s Server) handleErr(ctx context.Context, err error) {
	s.log.ErrorContext(ctx, err.Error())
}

func New(slogHandler slog.Handler, db database.DatabaseFace) Server {
	return Server{
		log: slog.New(slogHandler).With("service", "server"),
		db:  db,
	}
}
