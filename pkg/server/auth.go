package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s Server) auth() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", s.login())

	return r
}

func (s Server) login() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		req := LoginReq{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return http.StatusBadRequest, nil, err
		}

		fmt.Println(req)

		token, expiresIn, err := s.authPkg.CreateLoginToken(ctx, req.Email, req.Password, req.LongLivedToken)
		if err != nil {
			return http.StatusUnauthorized, nil, ErrUnauthorized{}
		}

		resp := LoginResp{
			Token:     token,
			TokenType: "Bearer",
			ExpiresIn: expiresIn,
		}

		return http.StatusOK, resp, nil
	},
	)
}
