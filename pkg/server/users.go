package server

import (
	"net/http"

	"github.com/bragemusic/core/pkg/internalusers"
)

func (s *Server) listUsers() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		users, err := s.authPkg.ListUsers(ctx)
		if err != nil {
			return Response{}, err
		}

		users = append(users, internalusers.GetIntenalUsers()...)

		return Response{Status: http.StatusOK, Payload: users}, nil
	})
}
