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

func (s Server) addPlaylist() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()

		plist := types.Playlist{}
		if err := json.NewDecoder(r.Body).Decode(&plist); err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		if err := s.mediamgr.AddPlaylist(ctx, plist); err != nil {
			return utils.Ptr(http.StatusInternalServerError), err
		}

		return utils.Ptr(http.StatusCreated), nil
	})
}

func (s Server) getPlaylist() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		plistID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		plist, err := s.mediamgr.GetPlaylist(ctx, plistID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "playlistID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, plist, nil
	},
	)
}
