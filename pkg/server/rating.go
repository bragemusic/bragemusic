package server

import (
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func (s *Server) getRating() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		ratingID, err := getParameter[uuid.UUID](ctx, "ratingID")
		if err != nil {
			return Response{}, err
		}

		rating, err := s.mediamgr.GetRating(ctx, ratingID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: rating}, nil
	})
}
