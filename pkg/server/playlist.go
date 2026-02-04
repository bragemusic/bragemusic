package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
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

func (s Server) addPlaylistTrack() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()

		plistID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return utils.Ptr(http.StatusForbidden), err
		}

		pt := PlaylistTrackReq{}
		if err := json.NewDecoder(r.Body).Decode(&pt); err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		if err := s.mediamgr.AddPlaylistTrack(ctx, plistID, pt.AlbumID, pt.TrackID, user.ID); err != nil {
			return utils.Ptr(http.StatusInternalServerError), err
		}

		return utils.Ptr(http.StatusCreated), nil
	})
}

func (s Server) deletePlaylist() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()

		pID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return utils.Ptr(int(http.StatusForbidden)), err
		}

		err = s.mediamgr.DeletePlaylist(ctx, pID, user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return utils.Ptr(http.StatusBadRequest), ErrIDNotFound{
					idKey: "playlistID",
					err:   err,
				}
			} else {
				return utils.Ptr(http.StatusInternalServerError), err
			}
		}

		return utils.Ptr(http.StatusNoContent), nil
	},
	)
}

func (s Server) getPlaylist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		plistID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return Response{}, err
		}

		plist, err := s.mediamgr.GetPlaylist(ctx, plistID)
		if err != nil {
			return Response{}, err
		}

		return Response{
			Payload: plist,
			Status:  http.StatusOK,
		}, nil
	},
	)
}

func (s Server) updatePlaylist() http.HandlerFunc {
	return s.handleVoid(func(w http.ResponseWriter, r *http.Request) (*int, error) {
		ctx := r.Context()

		id, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return utils.Ptr(http.StatusForbidden), err
		}

		plist := types.Playlist{}
		if err := json.NewDecoder(r.Body).Decode(&plist); err != nil {
			return utils.Ptr(http.StatusBadRequest), err
		}

		if err := s.mediamgr.UpdatePlaylist(ctx, id, plist, user.ID); err != nil {
			return utils.Ptr(http.StatusInternalServerError), err
		}

		return utils.Ptr(http.StatusNoContent), nil
	})
}
