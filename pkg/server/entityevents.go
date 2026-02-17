package server

import (
	"net/http"
)

func (s *Server) getEntityEvents() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		events, err := s.mediamgr.ListEntityEvents(ctx)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: events}, nil
	})
}
