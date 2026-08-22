package server

import (
	"net/http"

	"github.com/bragemusic/bragemusic/internal/vars"
	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) buildMount(ro []routes.RouteHandler) http.Handler {
	r := chi.NewRouter()

	for _, route := range ro {
		mws := []func(http.Handler) http.Handler{}

		for _, mw := range route.Middlewares() {
			mws = append(mws, mw.Func)
		}

		if len(route.Roles()) > 0 {
			mws = append(mws, s.authPkg.RoleCheckMiddleware(route.Roles()...))
			r.With(s.authPkg.RoleCheckMiddleware(route.Roles()...)).Method(route.Method(), route.Path(), route.Handler(s.log, s.errLog, &s.berr))
		}

		r.With(mws...).Method(route.Method(), route.Path(), route.Handler(s.log, s.errLog, &s.berr))
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

	r.Mount("/admin", s.buildMount(s.adminRoutes()))
	r.Mount("/albums", s.buildMount(s.albumRoutes()))
	r.Mount("/album-artists", s.buildMount(s.albumArtistRoutes()))
	r.Mount("/album-tracks", s.buildMount(s.albumTrackRoutes()))
	r.Mount("/artists", s.buildMount(s.artistRoutes()))
	r.Mount("/import", s.buildMount(s.importRoutes()))
	r.Mount("/likes", s.buildMount(s.likeRoutes()))
	r.Mount("/mediafiles", s.buildMount(s.mediafileRoutes()))
	r.Mount("/playlists", s.buildMount(s.playlistRoutes()))
	r.Mount("/playlist-tracks", s.buildMount(s.playlistTrackRoutes()))
	r.Mount("/ratings", s.buildMount(s.ratingRoutes()))
	r.Mount("/sync", s.buildMount(s.syncRoutes()))
	r.Mount("/search", s.buildMount(s.searchRoutes()))
	r.Mount("/smart-playlists", s.buildMount(s.smartPlaylistRoutes()))
	r.Mount("/tracks", s.buildMount(s.trackRoutes()))
	r.Mount("/track-analysis", s.buildMount(s.trackAnalysisRoutes()))
	r.Mount("/track-artists", s.buildMount(s.trackArtistRoutes()))

	r.Mount("/devices", s.buildMount(s.deviceRoutes()))
	r.Mount("/users", s.buildMount(s.userRoutes()))

	return r
}

func (s *Server) apiInfo() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		return Response{
			Payload: types.ServerApiInfo{
				Version: vars.VERSION,
				ServerInfo: types.ServerInfo{
					Name:        s.config.General.Name,
					Application: vars.SERVER_APP_NAME,
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
