package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
)

func (s Server) getMediaFile() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		mediafileID, err := getParameter[uuid.UUID](ctx, "mediafileID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		track, err := s.mediamgr.GetMediaFile(ctx, mediafileID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "trackID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, track, nil
	},
	)
}

func (s Server) getMediaFileFile() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()
		mediafileID, err := getParameter[uuid.UUID](ctx, "mediafileID")
		if err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		// FIXME: Dont hardcode flac
		w.Header().Add("Content-Type", "audio/flac")

		err = s.mediamgr.GetMediaFileFile(ctx, mediafileID, w)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return utils.Ptr(http.StatusBadRequest), ErrIDNotFound{
					idKey: "trackID",
					err:   err,
				}
			} else {
				return utils.Ptr(http.StatusInternalServerError), err
			}
		}

		return nil, nil
	},
	)
}
