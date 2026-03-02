package server

import (
	"net/http"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) buildMount(ro []routes.RouteHandler) http.Handler {
	r := chi.NewRouter()

	for _, route := range ro {
		if len(route.Roles()) > 0 {
			r.With(s.authPkg.RoleCheckMiddleware(route.Roles()...)).Method(route.Method(), route.Path(), route.Handler(s.log, s.errLog, &s.berr))
		} else {
			r.Method(route.Method(), route.Path(), route.Handler(s.log, s.errLog, &s.berr))
		}
	}

	return r
}

func (s *Server) api() http.Handler {
	r := chi.NewRouter()

	r.Use(s.authPkg.Middleware)

	r.Get("/info", s.apiInfo())
	r.Get("/user", s.user())

	r.Get("/img/*", s.getImage())
	r.Post("/img/artists/{artistID}", s.addImage(ArtistImage))
	r.Post("/img/albums/{albumID}", s.addImage(AlbumImage))
	r.Post("/img/playlists/{playlistID}", s.addImage(PlaylistImage))

	r.Mount("/albums", s.buildMount(s.albumRoutes()))
	r.Mount("/album-artists", s.buildMount(s.albumArtistRoutes()))
	r.Mount("/album-tracks", s.buildMount(s.albumTrackRoutes()))
	r.Mount("/artists", s.buildMount(s.artistRoutes()))
	r.Mount("/mediafiles", s.buildMount(s.mediafileRoutes()))
	r.Mount("/playlists", s.buildMount(s.playlistRoutes()))
	r.Mount("/playlist-tracks", s.buildMount(s.playlistTrackRoutes()))
	r.Mount("/sync", s.buildMount(s.syncRoutes()))
	r.Mount("/tracks", s.buildMount(s.trackRoutes()))

	r.Get("/ratings/{ratingID}", s.getRating())

	r.Get("/search", s.search())

	r.With(s.authPkg.RoleCheckMiddleware(types.UserRoleAdmin, types.UserRoleImporterWrite)).Post("/import/album", s.importAlbum())

	r.With(s.authPkg.RoleCheckMiddleware(types.UserRoleAdmin)).Get("/admin/entity-events", s.getEntityEvents())

	r.With(s.authPkg.RoleCheckMiddleware(types.UserRoleAdmin, types.UserRoleUsersGet)).Get("/users", s.listUsers())

	return r
}

func (s *Server) apiInfo() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		return Response{
			Payload: types.ServerApiInfo{
				Name:    s.config.Name,
				Version: config.VERSION,
				ServerInfo: types.ServerInfo{
					Application: config.SERVER_APP_NAME,
					Status:      types.HealthzRunning,
					ID:          uuid.Nil,
				},
			},
			Status: http.StatusOK,
		}, nil
	})
}

func (s *Server) user() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()
		user, err := s.authPkg.GetUserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		return Response{
			Payload: user,
			Status:  http.StatusOK,
		}, err
	})
}
