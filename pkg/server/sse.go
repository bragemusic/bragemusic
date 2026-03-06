package server

import (
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/go-chi/chi/v5"
)

func (s *Server) sse() http.Handler {
	r := chi.NewRouter()

	r.Use(s.authPkg.Middleware)

	r.Mount("/", s.buildMount(s.sseRoutes()))

	return r
}

func (s *Server) sseRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("POST", "/subscribe", s.sseHub.Handler(), nil, routes.RouteMeta{
			Summary:             "Subscribe to SSE events.",
			Description:         "Streams events from the server to the client.",
			ExpectedDescription: "Streamed event data",
			Tags:                []string{},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}
