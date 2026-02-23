package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) getArtist() http.HandlerFunc {
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

func (s *Server) listArtists() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		cnt, err := s.mediamgr.CountArtists(ctx)
		if err != nil {
			return Response{}, err
		}

		if r.URL.Query().Get("count") == "true" {
			return Response{Status: http.StatusOK, Payload: ListPayload[types.ArtistDetailed]{Count: cnt}}, nil
		}

		artists, err := s.mediamgr.ListArtists(ctx, database.SortByName, database.SortAsc)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: ListPayload[types.ArtistDetailed]{Count: cnt, Items: artists}}, nil
	})
}

func (s *Server) updateArtist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		artistID, err := getParameter[uuid.UUID](ctx, "artistID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		artist := types.Artist{}
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.UpdateArtist(ctx, artistID, artist, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
