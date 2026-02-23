package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) addPlayHistory() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		plist := types.Playlist{}
		if err := json.NewDecoder(r.Body).Decode(&plist); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.AddPlayCount(ctx, trackID, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusCreated}, nil
	})
}
