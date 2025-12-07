package server

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s Server) getImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "*")
		fp := strings.TrimPrefix(r.URL.Path, pathPrefix)

		filename := filepath.Join(s.config.Paths.ImageDir, fp)

		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, filename)
	}
}
