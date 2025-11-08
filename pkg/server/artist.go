package server

import "net/http"

func (s Server) listArtists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		artists, err := s.db.ListArtists(ctx)
		if err != nil {
			s.handleErr(ctx, err)
			return
		}

		if err = s.writeJSON(ctx, w, artists); err != nil {
			s.handleErr(ctx, err)
			return
		}
	}
}
