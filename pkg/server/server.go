package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/go-chi/chi/v5"
)

type (
	handlerFunc func(http.ResponseWriter, *http.Request) (Response, error)
)

type Server struct {
	log      *slog.Logger
	errLog   *slog.Logger
	mediamgr *mediamanager.MediaManager
	authPkg  *auth.Auth
	importer *importer.Importer
	jobmgr   *jobmanager.JobManager
	config   Config
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware(*s.log, []string{"/healthz"}))

	r.Get("/healthz", s.healthz())

	r.Mount("/api", s.api())
	r.Mount("/auth", s.auth())

	return r
}

func (s Server) healthz() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		return Response{Status: http.StatusOK, Payload: Status{
			Application: "brage-server", // hardcoded
			Name:        "Brage Server", // from config
			Version:     "v0.0.1",
			Status:      HealthzRunning,
		}}, nil
	})
}

func (s Server) handle(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := f(w, r)
		ctx := r.Context()
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, s.errLog)
			return
		}

		if resp.Payload != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(resp.Status)

			err = json.NewEncoder(w).Encode(resp.Payload)
			if err != nil {
				s.log.ErrorContext(ctx, err.Error())
			}
		} else {
			w.WriteHeader(resp.Status)
		}
	}
}

func (s Server) Start(ctx context.Context) error {
	go s.jobmgr.StartScheduler(ctx)

	s.log.InfoContext(ctx, fmt.Sprintf("serving on port %d", s.config.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.Handler()); err != nil {
		return err
	}

	return nil
}

func New(slogHandler slog.Handler, m *mediamanager.MediaManager, a *auth.Auth, i *importer.Importer, j *jobmanager.JobManager, c Config) Server {
	return Server{
		log:      slog.New(slogHandler).With("service", "server"),
		errLog:   slog.New(slogHandler),
		mediamgr: m,
		config:   c,
		authPkg:  a,
		importer: i,
		jobmgr:   j,
	}
}
