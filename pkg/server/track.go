package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) addTrackRating() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return Response{}, err
		}

		rating := types.RatingReq{}
		if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.RateTrack(ctx, trackID, user.ID, rating.Value); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}

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

func (s *Server) getTrackDetailed() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		track, err := s.mediamgr.GetTrackDetailed(ctx, trackID, albumID, user.ID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: track}, nil
	},
	)
}

func (s *Server) getTrackRatings() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
		if err != nil {
			return Response{}, err
		}

		ratings, err := s.mediamgr.GetTrackRatings(ctx, trackID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: ratings}, nil
	})
}

func (s *Server) listTracks() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		cnt, err := s.mediamgr.CountTracks(ctx)
		if err != nil {
			return Response{}, err
		}

		if r.URL.Query().Get("count") == "true" {
			return Response{Status: http.StatusOK, Payload: types.ListPayload[types.TrackDetailed]{Count: cnt}}, nil
		}

		return Response{Status: http.StatusOK, Payload: types.ListPayload[types.TrackDetailed]{Count: cnt}}, nil
	})
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

func (s *Server) listAlbumTracksDetailed() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		tracks, err := s.mediamgr.ListTracksDetailedByAlbum(ctx, albumID, user.ID)
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

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		track := types.TrackUpdate{}
		if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.UpdateTrack(ctx, id, track, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
