package server

import (
	"net/http"

	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) search() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		searchTerm := r.URL.Query().Get("q")

		searchItems, err := s.mediamgr.SearchFull(ctx, searchTerm)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: types.ListPayload[types.SearchItem]{Count: len(searchItems), Items: searchItems}}, nil
	})
}
