package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) searchRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/media", s.search(), nil, routes.RouteMeta{
			Summary:             "Search for media entites",
			Description:         "Search for artists, albums, tracks and playlists. Returns a list with ids and names.",
			ExpectedDescription: "Search results",
			Tags:                []string{"Search"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) search() routes.RouteFunc[ReqSearch, types.ListPayload[types.SearchItem]] {
	return func(ctx context.Context, req ReqSearch, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.SearchItem]], err error) {
		searchItems, err := s.mediamgr.SearchFull(ctx, req.SearchTerm)
		if err != nil {
			return resp, err
		}

		return types.Response[types.ListPayload[types.SearchItem]]{
			Payload: types.ListPayload[types.SearchItem]{
				Items: searchItems,
				Count: len(searchItems),
			},
			Status: http.StatusOK,
		}, nil
	}
}
