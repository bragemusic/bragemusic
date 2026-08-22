package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) trackAnalysisRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{trackAnalysisID}", s.getTrackAnalysisByID(), nil, routes.RouteMeta{
			Summary:             "Retrieve a track analysis by ID.",
			Description:         "Returns metadata about the specified track analysis.",
			ExpectedDescription: "Metadata about the track analysis",
			Tags:                []string{"Track Analysis"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getTrackAnalysisByID() routes.RouteFunc[ReqTrackAnalysisGet, types.TrackAnalysis] {
	return func(ctx context.Context, req ReqTrackAnalysisGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.TrackAnalysis], err error) {
		analysis, err := s.mediamgr.GetTrackAnalysisByID(ctx, req.ID)
		if err != nil {
			return types.Response[types.TrackAnalysis]{}, err
		}

		return types.Response[types.TrackAnalysis]{
			Payload: analysis,
			Status:  http.StatusOK,
		}, nil
	}
}
