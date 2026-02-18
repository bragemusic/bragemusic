package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/sync/errgroup"
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
	httpSrv  *http.Server
	ready    atomic.Bool
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware(*s.log, []string{"/healthz"}))

	r.Get("/healthz", s.healthz())
	r.Get("/readyz", s.readyz())
	r.Get("/info", s.info())

	r.Mount("/api", s.api())
	r.Mount("/auth", s.auth())

	return r
}

func (s *Server) healthz() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		return Response{Status: http.StatusOK, Payload: Healthz{
			Status: HealthzRunning,
		}}, nil
	})
}

func (s *Server) info() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		return Response{Status: http.StatusOK, Payload: ServerInfo{
			Application: config.SERVER_APP_NAME,
			ID:          uuid.Nil,
			Status:      HealthzRunning,
		}}, nil
	})
}

func (s *Server) readyz() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		if !s.ready.Load() {
			return Response{Status: http.StatusServiceUnavailable, Payload: nil}, nil
		}
		return Response{Status: http.StatusServiceUnavailable, Payload: nil}, nil
	})
}

func (s *Server) handle(f handlerFunc) http.HandlerFunc {
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

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.Handler(),
	}

	// Mark as ready
	s.ready.Store(true)

	g.Go(func() error {
		s.log.InfoContext(ctx, "HTTP server starting", "port", s.config.Port)

		if err := s.httpSrv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		s.jobmgr.StartScheduler(ctx)
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		s.log.InfoContext(ctx, "shutdown initiated")

		// Fail readiness immediately
		s.ready.Store(false)

		// Give Kubernetes time to stop routing traffic
		time.Sleep(2 * time.Second)

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			25*time.Second, // must be < k8s terminationGracePeriodSeconds
		)
		defer shutdownCancel()

		return s.httpSrv.Shutdown(shutdownCtx)
	})

	return g.Wait()
}

// func (s Server) Start(ctx context.Context) error {
// 	go s.jobmgr.StartScheduler(ctx)

// 	s.log.InfoContext(ctx, fmt.Sprintf("serving on port %d", s.config.Port))
// 	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.Handler()); err != nil {
// 		return err
// 	}

// 	return nil
// }

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
