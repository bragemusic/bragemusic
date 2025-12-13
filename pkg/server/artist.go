package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/go-chi/chi/v5"
)

func (s Server) getArtist() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		artistID := chi.URLParamFromCtx(ctx, "artistID")
		if artistID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "artistID",
				err:   errors.New("could not parse artistID"),
			}
		}

		artist, err := s.mediamgr.GetArtist(ctx, artistID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "artistID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, artist, nil
	})
}

func (s Server) listArtists() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		artists, err := s.mediamgr.ListArtists(ctx, database.SortByName, database.SortAsc)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		return http.StatusOK, artists, nil
	})
}
