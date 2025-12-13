package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/go-chi/chi/v5"
)

func (s Server) getAlbum() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumID := chi.URLParamFromCtx(ctx, "albumID")
		if albumID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "albumID",
				err:   errors.New("could not parse albumID"),
			}
		}

		album, err := s.mediamgr.GetAlbum(ctx, albumID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, album, nil
	})
}

func (s Server) listAlbums() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		artistID := chi.URLParamFromCtx(ctx, "artistID")
		if artistID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "artistID",
				err:   errors.New("could not parse artistID"),
			}
		}

		albums, err := s.mediamgr.ListAlbumsByArtist(ctx, artistID, database.SortByDate, database.SortAsc)
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

		return http.StatusOK, albums, nil
	})
}
