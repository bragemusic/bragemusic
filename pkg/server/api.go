package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s Server) api() http.Handler {
	r := chi.NewRouter()

	r.Use(s.authPkg.Middleware)

	r.Get("/status", s.status())
	r.Get("/user", s.user())

	r.Get("/img/*", s.getImage())

	r.Get("/artists", s.listArtists())
	r.Get("/artists/{artistID}", s.getArtist())
	r.Get("/artists/{artistID}/albums", s.listAlbums())

	r.Get("/albums/{albumID}", s.getAlbum())
	r.Get("/albums/{albumID}/tracks", s.listAlbumTracks())

	r.Get("/tracks/{trackID}", s.getTrack())
	r.Get("/tracks/{trackID}/file", s.getTrackFile())

	r.Post("/sync", s.sync())
	r.Post("/sync/play-history", s.syncPlayHistory())

	return r
}

func (s Server) status() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		return http.StatusOK, Status{
			Application: "brage-server",
			Version:     "v0.0.1",
			Status:      HealthzRunning,
		}, nil
	})
}

func (s Server) user() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()
		user, err := s.authPkg.GetUserFromContext(ctx)
		if err != nil {
			return http.StatusForbidden, nil, err
		}

		return http.StatusOK, user, nil
	})
}
