package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s Server) getArtist() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		artistID, err := getParameter[uuid.UUID](ctx, "artistID")
		if err != nil {
			return http.StatusBadRequest, nil, err
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

func (s Server) updateArtist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		artistID, err := getParameter[uuid.UUID](ctx, "artistID")
		if err != nil {
			return Response{}, err
		}

		artist := types.Artist{}
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.UpdateArtist(ctx, artistID, artist); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
