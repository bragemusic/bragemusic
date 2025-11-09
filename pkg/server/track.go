package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func (s Server) getTrack() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		trackID := chi.URLParamFromCtx(ctx, "trackID")
		if trackID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "trackID",
				err:   errors.New("could not parse trackID"),
			}
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

func (s Server) getTrackFile() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()

		trackID := chi.URLParamFromCtx(ctx, "trackID")
		if trackID == "" {
			return utils.Ptr(http.StatusBadRequest), ErrIDNotFound{
				idKey: "trackID",
				err:   errors.New("could not parse trackID"),
			}
		}

		w.Header().Add("Content-Type", "audio/flac")

		err := s.mediamgr.GetTrackFile(ctx, trackID, w)
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

func (s Server) listAlbumTracks() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumID := chi.URLParamFromCtx(ctx, "albumID")
		if albumID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "albumID",
				err:   errors.New("could not parse albumID"),
			}
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
