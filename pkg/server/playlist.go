package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) addPlaylist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		plist := types.Playlist{}
		if err := json.NewDecoder(r.Body).Decode(&plist); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.AddPlaylist(ctx, plist, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusCreated}, nil
	})
}

func (s *Server) addPlaylistTrack() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		plistID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		pt := types.PlaylistTrackReq{}
		if err := json.NewDecoder(r.Body).Decode(&pt); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.AddPlaylistTrack(ctx, plistID, pt.AlbumID, pt.TrackID, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusCreated}, nil
	})
}

func (s *Server) deletePlaylist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		pID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		err = s.mediamgr.DeletePlaylist(ctx, pID, user.ID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	},
	)
}

func (s *Server) getPlaylist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		plistID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		plist, err := s.mediamgr.GetPlaylist(ctx, plistID, user.ID)
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

func (s *Server) listPlaylists() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		cnt, err := s.mediamgr.CountPlaylists(ctx, user.ID)
		if err != nil {
			return Response{}, err
		}

		if r.URL.Query().Get("count") == "true" {
			return Response{Status: http.StatusOK, Payload: types.ListPayload[types.Playlist]{Count: cnt}}, nil
		}

		sB := database.SortBy(r.URL.Query().Get("sortBy"))
		sO := database.SortOrder(r.URL.Query().Get("sortOrder"))
		iP := r.URL.Query().Get("includePublic") == "true"

		plists, err := s.mediamgr.ListPlaylists(ctx, user.ID, iP, sB, sO)
		if err != nil {
			return Response{}, err
		}

		return Response{
			Payload: types.ListPayload[types.Playlist]{Count: cnt, Items: plists},
			Status:  http.StatusOK,
		}, nil
	},
	)
}

func (s *Server) listPlaylistTracks() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		plistID, err := getParameter[uuid.UUID](ctx, "playlistID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		cnt, err := s.mediamgr.CountPlaylistTracks(ctx, plistID, user.ID)
		if err != nil {
			return Response{}, err
		}

		if r.URL.Query().Get("count") == "true" {
			return Response{Status: http.StatusOK, Payload: types.ListPayload[types.TrackDetailed]{Count: cnt}}, nil
		}

		sB := database.SortBy(r.URL.Query().Get("sortBy"))
		sO := database.SortOrder(r.URL.Query().Get("sortOrder"))

		tracks, err := s.mediamgr.ListPlaylistTracks(ctx, plistID, user.ID, sB, sO)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: types.ListPayload[types.TrackDetailed]{Count: cnt, Items: tracks}}, nil
	})
}

func (s *Server) updatePlaylist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		id, err := getParameter[uuid.UUID](ctx, "playlistID")
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

		if err := s.mediamgr.UpdatePlaylist(ctx, id, plist, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
