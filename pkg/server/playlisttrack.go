package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/gofrs/uuid/v5"
)

func (s Server) getPlaylistTrack() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		ptID, err := getParameter[uuid.UUID](ctx, "playlistTrackID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return http.StatusForbidden, nil, err
		}

		pt, err := s.mediamgr.GetPlaylistTrack(ctx, ptID, user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "playlistTrackID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, pt, nil
	},
	)
}
