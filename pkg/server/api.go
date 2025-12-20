package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s Server) api() http.Handler {
	r := chi.NewRouter()

	r.Use(s.authPkg.Middleware)

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
