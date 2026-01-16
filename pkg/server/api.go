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
	r.Post("/img/artists/{artistID}", s.addImage(ArtistImage))
	r.Post("/img/albums/{albumID}", s.addImage(AlbumImage))

	r.Get("/artists", s.listArtists())
	r.Get("/artists/{artistID}", s.getArtist())
	r.Put("/artists/{artistID}", s.updateArtist())
	r.Get("/artists/{artistID}/albums", s.listAlbums())

	r.Get("/albums/{albumID}", s.getAlbum())
	r.Put("/albums/{albumID}", s.updateAlbum())
	r.Get("/albums/{albumID}/tracks", s.listAlbumTracks())
	r.Get("/albums/{albumID}/artists/{artistID}/roles/{role}", s.getAlbumArtist())
	r.Get("/albums/{albumID}/disc/{discNumber}/track/{trackNumber}", s.getAlbumTrack())

	r.Get("/album-artists/{albumArtistID}", s.getAlbumArtistByID())

	r.Get("/tracks/{trackID}", s.getTrack())

	r.Get("/mediafiles/{mediafileID}", s.getMediaFile())
	r.Get("/mediafiles/{mediafileID}/file", s.getMediaFileFile())

	r.Post("/sync", s.sync())
	r.Post("/sync/play-history", s.syncPlayHistory())

	return r
}

func (s Server) status() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		return http.StatusOK, Status{
			Application: "brage-server", // hardcoded
			Name:        "Brage Server", // from config
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
