package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) smartPlaylistRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/content/{contentID}", s.getSmartPlaylistContent(), nil, routes.RouteMeta{
			Summary:             "Get smart content for a playlist by ID",
			Description:         "Returns metadata for the specified playlist smart content.",
			ExpectedDescription: "Playlist smart content",
			Tags:                []string{"Smart Playlist"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/artist/{artistID}", s.getSmartPlaylistArtist(), nil, routes.RouteMeta{
			Summary:             "Get smart artist for a playlist by ID",
			Description:         "Returns metadata for the specified playlist smart artist.",
			ExpectedDescription: "Playlist smart artist",
			Tags:                []string{"Smart Playlist"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getSmartPlaylistContent() routes.RouteFunc[ReqSmartPlaylistContentGet, types.SmartPlaylistContent] {
	return func(ctx context.Context, req ReqSmartPlaylistContentGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.SmartPlaylistContent], err error) {
		content, err := s.mediamgr.GetSmartPlaylistContent(ctx, req.ID, user.ID)
		if err != nil {
			return types.Response[types.SmartPlaylistContent]{}, err
		}

		return types.Response[types.SmartPlaylistContent]{
			Payload: content,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) getSmartPlaylistArtist() routes.RouteFunc[ReqSmartPlaylistArtistGet, types.SmartPlaylistArtist] {
	return func(ctx context.Context, req ReqSmartPlaylistArtistGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.SmartPlaylistArtist], err error) {
		artist, err := s.mediamgr.GetSmartPlaylistArtist(ctx, req.ID, user.ID)
		if err != nil {
			return types.Response[types.SmartPlaylistArtist]{}, err
		}

		return types.Response[types.SmartPlaylistArtist]{
			Payload: artist,
			Status:  http.StatusOK,
		}, nil
	}
}
