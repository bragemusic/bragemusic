package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
)

func (s Server) getTrack() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		track, err := s.mediamgr.GetTrack(ctx, trackID)
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

func (s Server) listAlbumTracks() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		tracks, err := s.mediamgr.ListTracksByAlbum(ctx, albumID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, tracks, nil
	})
}

func (s Server) updateTrack() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()

		id, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		track := types.TrackUpdate{}
		if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		if err := s.mediamgr.UpdateTrack(ctx, id, track); err != nil {
			return utils.Ptr(http.StatusInternalServerError), err
		}

		return utils.Ptr(http.StatusNoContent), nil
	})
}
