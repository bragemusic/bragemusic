package server

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(logger slog.Logger, excludePaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(ww, r)

			level := slog.LevelInfo
			if ww.Status() >= 400 {
				level = slog.LevelWarn
			} else if ww.Status() >= 500 {
				level = slog.LevelError
			} else if slices.Contains(excludePaths, r.URL.Path) {
				level = slog.LevelDebug
			}

			logger.Log(
				ctx,
				level,
				"http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration", time.Since(start),
			)
		})
	}
}
