package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/go-chi/chi/v5"
)

type (
	handlerFuncErrJson func(http.ResponseWriter, *http.Request) (int, any, error)
	handlerFuncErrVoid func(http.ResponseWriter, *http.Request) (*int, error)
)

type Server struct {
	log      *slog.Logger
	mediamgr *mediamanager.MediaManager
	auth     *auth.Auth
	config   Config
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", s.healthz())

	r.Mount("/api", s.api())

	return r
}

func (s Server) healthz() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		return http.StatusOK, Healthz{
			Application: "brage-server",
			Version:     "v0.0.1",
			Status:      HealthzRunning,
		}, nil
	})
}

func (s Server) handleVoid(f handlerFuncErrVoid) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := f(w, r)

		ctx := r.Context()
		if status != nil {
			w.WriteHeader(*status)
		}

		if err != nil {
			sErr, ok := err.(ServerError)
			if ok {
				jErr := json.NewEncoder(w).Encode(map[string]string{"error": sErr.UserError()})
				if jErr != nil {
					s.log.ErrorContext(ctx, err.Error())
				}
			} else {
				jErr := json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
				if jErr != nil {
					s.log.ErrorContext(ctx, err.Error())
				}
			}
			s.log.ErrorContext(ctx, err.Error())
			return
		}
	}
}

func (s Server) handleJSON(f handlerFuncErrJson) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, payload, err := f(w, r)

		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)

		if err != nil {
			sErr, ok := err.(ServerError)
			if ok {
				jErr := json.NewEncoder(w).Encode(map[string]string{"error": sErr.UserError()})
				if jErr != nil {
					s.log.ErrorContext(ctx, err.Error())
				}
			} else {
				jErr := json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
				if jErr != nil {
					s.log.ErrorContext(ctx, err.Error())
				}
			}
			s.log.ErrorContext(ctx, err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(payload)
		if err != nil {
			s.log.ErrorContext(ctx, err.Error())
		}
	}
}

func New(slogHandler slog.Handler, m *mediamanager.MediaManager, a *auth.Auth, c Config) Server {
	return Server{
		log:      slog.New(slogHandler).With("service", "server"),
		mediamgr: m,
		config:   c,
		auth:     a,
	}
}
