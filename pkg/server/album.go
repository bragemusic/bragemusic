package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s Server) listAlbums() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		artistID := chi.URLParamFromCtx(ctx, "artistID")
		if artistID == "" {
			s.handleErr(ctx, errors.New("no artistID provided"))
			return
		}

		albums, err := s.db.ListAlbumsByArtist(ctx, artistID)
		if err != nil {
			s.handleErr(ctx, err)
			return
		}

		if err = s.writeJSON(ctx, w, albums); err != nil {
			s.handleErr(ctx, err)
			return
		}
	}
}
