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
		routes.New("POST", "/search", s.listTracksDetailed(), nil, routes.RouteMeta{
			Summary:             "List all detailed tracks",
			Description:         "Detailed metadata of tracks",
			ExpectedDescription: "Tracks",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/liked", s.listLikedTracks(), nil, routes.RouteMeta{
			Summary:             "List liked tracks.",
			Description:         "Returns metadata about all tracks liked by the authenticated user, ordered by when they were liked.",
			ExpectedDescription: "Metadata about the liked tracks",
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
		routes.New("POST", "/{trackID}/like", s.addTrackLike(), nil, routes.RouteMeta{
			Summary:             "Like a track by ID",
			Description:         "Add a like to the selected track for the authed user.",
			ExpectedDescription: "Like accepted",
			Tags:                []string{"Tracks"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("DELETE", "/{trackID}/like", s.deleteTrackLike(), nil, routes.RouteMeta{
			Summary:             "Remove a track like by ID",
			Description:         "Remove a like to the selected track for the authed user.",
			ExpectedDescription: "Like removed",
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

func (s *Server) addTrackLike() routes.RouteFunc[ReqTracksGet, types.NoResponse] {
	return func(ctx context.Context, req ReqTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.mediamgr.AddLike(ctx, req.TrackID, user.ID)
		if err != nil {
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

func (s *Server) deleteTrackLike() routes.RouteFunc[ReqTracksGet, types.NoResponse] {
	return func(ctx context.Context, req ReqTracksGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.mediamgr.RemoveLike(ctx, req.TrackID, user.ID)
		if err != nil {
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

func (s *Server) listLikedTracks() routes.RouteFunc[ReqList, types.ListPayload[types.TrackDetailed]] {
	return func(ctx context.Context, req ReqList, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.TrackDetailed]], err error) {
		tracks, err := s.mediamgr.ListLikedTracksDetailed(ctx, user.ID)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.TrackDetailed]]{
				Payload: types.ListPayload[types.TrackDetailed]{
					Items: nil,
					Count: len(tracks),
				},
				Status: http.StatusOK,
			}, nil
		}

		return types.Response[types.ListPayload[types.TrackDetailed]]{
			Payload: types.ListPayload[types.TrackDetailed]{
				Items: tracks,
				Count: len(tracks),
			},
			Status: http.StatusOK,
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

func (s *Server) listTracksDetailed() routes.RouteFunc[ReqListTrackPagination, types.ListPaginationPayload[types.TrackDetailed]] {
	return func(ctx context.Context, req ReqListTrackPagination, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPaginationPayload[types.TrackDetailed]], err error) {
		items, page, limit, totalPages, totalItems, err := s.mediamgr.ListTracksDetailed(ctx, req.TrackFilter, req.Page, req.Limit)
		if err != nil {
			return types.Response[types.ListPaginationPayload[types.TrackDetailed]]{}, err
		}

		return types.Response[types.ListPaginationPayload[types.TrackDetailed]]{
			Payload: types.ListPaginationPayload[types.TrackDetailed]{
				Items:      items,
				Page:       page,
				Limit:      limit,
				TotalPages: totalPages,
				TotalItems: totalItems,
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
