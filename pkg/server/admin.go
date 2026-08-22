package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) adminRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/entity-events", s.getEntityEvents(), []types.UserRole{types.UserRoleAdmin}, routes.RouteMeta{
			Summary:             "List all entity events.",
			Description:         "Return entity events that has occured on the server. An entity event is a CRUD operation on the database.",
			ExpectedDescription: "List of events",
			Tags:                []string{"Admin"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getEntityEvents() routes.RouteFunc[ReqList, types.ListPayload[types.EntityEvent]] {
	return func(ctx context.Context, req ReqList, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.EntityEvent]], err error) {
		events, err := s.mediamgr.ListEntityEvents(ctx)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.EntityEvent]]{
				Payload: types.ListPayload[types.EntityEvent]{
					Items: nil,
					Count: len(events),
				},
				Status: http.StatusOK,
			}, nil
		}

		return types.Response[types.ListPayload[types.EntityEvent]]{
			Payload: types.ListPayload[types.EntityEvent]{
				Items: events,
				Count: len(events),
			},
			Status: http.StatusOK,
		}, nil
	}
}
