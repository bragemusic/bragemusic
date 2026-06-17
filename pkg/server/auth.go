package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/types"
	"github.com/go-chi/chi/v5"
)

func (s *Server) auth() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", s.login())
	r.With(s.authPkg.Middleware).Get("/logout", s.logout())

	r.With(s.authPkg.Middleware).Mount("/tokens", s.buildMount(s.tokenRoutes()))
	r.With(s.authPkg.Middleware).Mount("/user-roles", s.buildMount(s.userRoleRoutes()))

	return r
}

func (s *Server) logout() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()
		http.SetCookie(w, &http.Cookie{
			Name:    "brage_session_token",
			Path:    "/",
			Value:   "",
			Expires: time.Unix(0, 0),
		})

		tokenID, err := auth.TokenIDFromContext(ctx)
		if err != nil {
			s.log.ErrorContext(ctx, "could not get token from context to remove token on logout", "error", err.Error())
			return Response{Status: http.StatusNoContent, Payload: nil}, nil
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			s.log.ErrorContext(ctx, "could not get userID from context to remove token on logout", "error", err.Error())
			return Response{Status: http.StatusNoContent, Payload: nil}, nil
		}

		err = s.authPkg.RemoveToken(ctx, tokenID, user.ID)
		if err != nil {
			s.log.ErrorContext(ctx, "could not remove token", "user_id", user.ID, "token_id", tokenID, "error", err.Error())
			return Response{Status: http.StatusNoContent, Payload: nil}, nil
		}

		return Response{Status: http.StatusNoContent, Payload: nil}, nil
	},
	)
}

func (s *Server) login() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		req := types.LoginReq{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return Response{}, err
		}

		token, expiresIn, err := s.authPkg.CreateLoginToken(ctx, req.Email, req.Password, req.LongLivedToken)
		if err != nil {
			return Response{}, err
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "brage_session_token",
			Path:    "/",
			Value:   token,
			Expires: time.Now().Add(time.Duration(expiresIn) * time.Second),
		})

		resp := types.LoginResp{
			Token:     token,
			TokenType: "Bearer",
			ExpiresIn: expiresIn,
		}

		return Response{Status: http.StatusOK, Payload: resp}, nil
	},
	)
}
