package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) getTrack() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return Response{}, err
		}

		track, err := s.mediamgr.GetTrack(ctx, trackID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: track}, nil
	},
	)
}

func (s *Server) listAlbumTracks() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		tracks, err := s.mediamgr.ListTracksByAlbum(ctx, albumID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: tracks}, nil
	})
}

func (s *Server) updateTrack() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		id, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return Response{}, err
		}

		track := types.TrackUpdate{}
		if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.UpdateTrack(ctx, id, track); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
