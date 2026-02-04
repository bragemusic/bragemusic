package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s Server) getArtist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		artistID, err := getParameter[uuid.UUID](ctx, "artistID")
		if err != nil {
			return Response{}, err
		}

		artist, err := s.mediamgr.GetArtist(ctx, artistID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: artist}, nil
	})
}

func (s Server) listArtists() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		artists, err := s.mediamgr.ListArtists(ctx, database.SortByName, database.SortAsc)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: artists}, nil
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
