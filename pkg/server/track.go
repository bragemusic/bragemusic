package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) trackRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listTracks(), nil, routes.RouteMeta{
			Summary:             "List all tracks.",
			Description:         "Returns metadata about all tracks.",
			ExpectedDescription: "Metadata about the tracks",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{trackID}", s.getTrack(), nil, routes.RouteMeta{
			Summary:             "Retrieve a track by ID.",
			Description:         "Returns metadata about the specified track.",
			ExpectedDescription: "Metadata about the track",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("PUT", "/{trackID}", s.updateTrack(), nil, routes.RouteMeta{
			Summary:             "Update a track by ID.",
			Description:         "Updates metadata about the specified track.",
			ExpectedDescription: "Update succeded",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("POST", "/{trackID}/play-history", s.addPlayHistory(), nil, routes.RouteMeta{
			Summary:             "Add play history entry to a track",
			Description:         "Adds one play history event for a track.",
			ExpectedDescription: "Play history accepted",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("GET", "/{trackID}/ratings", s.getTrackRatings(), nil, routes.RouteMeta{
			Summary:             "List all ratings of a track by ID.",
			Description:         "Lists all ratings from all users of a track.",
			ExpectedDescription: "List of ratings",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("POST", "/{trackID}/ratings", s.addTrackRating(), nil, routes.RouteMeta{
			Summary:             "Rate a track by ID.",
			Description:         "Add or update a rating for the track with to the logged in user",
			ExpectedDescription: "Rating accepted",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
	}
}

func (s *Server) addPlayHistory() routes.RouteFunc[ReqTracksGet, types.NoResponse] {
	return func(ctx context.Context, req ReqTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.AddPlayCount(ctx, req.TrackID, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) addTrackRating() routes.RouteFunc[ReqTracksAddRating, types.NoResponse] {
	return func(ctx context.Context, req ReqTracksAddRating, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.RateTrack(ctx, req.TrackID, user.ID, req.RatingReq.Value); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) getTrack() routes.RouteFunc[ReqTracksGet, types.Track] {
	return func(ctx context.Context, req ReqTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.Track], err error) {
		track, err := s.mediamgr.GetTrack(ctx, req.TrackID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.Track]{
			Payload: track,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) getTrackRatings() routes.RouteFunc[ReqTracksGet, []types.Rating] {
	return func(ctx context.Context, req ReqTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[[]types.Rating], err error) {
		ratings, err := s.mediamgr.GetTrackRatings(ctx, req.TrackID)
		if err != nil {
			return resp, err
		}

		return types.Response[[]types.Rating]{
			Payload: ratings,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) listTracks() routes.RouteFunc[ReqList, types.ListPayload[types.TrackDetailed]] {
	return func(ctx context.Context, req ReqList, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.TrackDetailed]], err error) {
		cnt, err := s.mediamgr.CountTracks(ctx)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.TrackDetailed]]{
				Payload: types.ListPayload[types.TrackDetailed]{
					Items: nil,
					Count: cnt,
				},
				Status: http.StatusOK,
			}, nil
		}

		return types.Response[types.ListPayload[types.TrackDetailed]]{
			Payload: types.ListPayload[types.TrackDetailed]{
				Items: nil,
				Count: cnt,
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) updateTrack() routes.RouteFunc[ReqTracksUpdate, types.NoResponse] {
	return func(ctx context.Context, req ReqTracksUpdate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.UpdateTrack(ctx, req.TrackID, req.TrackUpdate, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

///

// func (s *Server) addTrackRating() http.HandlerFunc {
// 	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
// 		ctx := r.Context()

// 		user, err := auth.UserFromContext(ctx)
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		rating := types.RatingReq{}
// 		if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
// 			return Response{}, err
// 		}

// 		if err := s.mediamgr.RateTrack(ctx, trackID, user.ID, rating.Value); err != nil {
// 			return Response{}, err
// 		}

// 		return Response{Status: http.StatusNoContent}, nil
// 	})
// }

// func (s *Server) getTrack() http.HandlerFunc {
// 	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
// 		ctx := r.Context()

// 		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		track, err := s.mediamgr.GetTrack(ctx, trackID)
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		return Response{Status: http.StatusOK, Payload: track}, nil
// 	},
// 	)
// }

// func (s *Server) getTrackRatings() http.HandlerFunc {
// 	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
// 		ctx := r.Context()

// 		trackID, err := getParameter[uuid.UUID](ctx, "trackID")
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		ratings, err := s.mediamgr.GetTrackRatings(ctx, trackID)
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		return Response{Status: http.StatusOK, Payload: ratings}, nil
// 	})
// }

// func (s *Server) listTracks() http.HandlerFunc {
// 	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
// 		ctx := r.Context()

// 		cnt, err := s.mediamgr.CountTracks(ctx)
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		if r.URL.Query().Get("count") == "true" {
// 			return Response{Status: http.StatusOK, Payload: types.ListPayload[types.TrackDetailed]{Count: cnt}}, nil
// 		}

// 		return Response{Status: http.StatusOK, Payload: types.ListPayload[types.TrackDetailed]{Count: cnt}}, nil
// 	})
// }

// func (s *Server) updateTrack() http.HandlerFunc {
// 	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
// 		ctx := r.Context()

// 		id, err := getParameter[uuid.UUID](ctx, "trackID")
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		user, err := auth.UserFromContext(ctx)
// 		if err != nil {
// 			return Response{}, err
// 		}

// 		track := types.TrackUpdate{}
// 		if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
// 			return Response{}, err
// 		}

// 		if err := s.mediamgr.UpdateTrack(ctx, id, track, user.ID); err != nil {
// 			return Response{}, err
// 		}

// 		return Response{Status: http.StatusNoContent}, nil
// 	})
// }
