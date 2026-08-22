package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) albumArtistRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{albumArtistID}", s.getAlbumArtistByID(), nil, routes.RouteMeta{
			Summary:             "Retrieve an album artist by ID.",
			Description:         "Returns metadata about the specified album artist.",
			ExpectedDescription: "Metadata about the album artist",
			Tags:                []string{"Album Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) getAlbumArtistByID() routes.RouteFunc[ReqAlbumArtistsGet, types.AlbumArtist] {
	return func(ctx context.Context, req ReqAlbumArtistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.AlbumArtist], err error) {
		albumArtist, err := s.mediamgr.GetAlbumArtistByID(ctx, req.AlbumArtistID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.AlbumArtist]{
			Payload: albumArtist,
			Status:  http.StatusOK,
		}, nil
	}
}
