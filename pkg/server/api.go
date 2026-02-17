package server

import (
	"net/http"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) api() http.Handler {
	r := chi.NewRouter()

	r.Use(s.authPkg.Middleware)

	r.Get("/info", s.apiInfo())
	r.Get("/user", s.user())

	r.Get("/img/*", s.getImage())
	r.Post("/img/artists/{artistID}", s.addImage(ArtistImage))
	r.Post("/img/albums/{albumID}", s.addImage(AlbumImage))
	r.Post("/img/playlists/{playlistID}", s.addImage(PlaylistImage))

	r.Get("/artists", s.listArtists())
	r.Get("/artists/{artistID}", s.getArtist())
	r.Put("/artists/{artistID}", s.updateArtist())
	r.Get("/artists/{artistID}/albums", s.listAlbums())

	r.Get("/albums/{albumID}", s.getAlbum())
	r.Put("/albums/{albumID}", s.updateAlbum())
	r.Get("/albums/{albumID}/tracks", s.listAlbumTracks())
	r.Get("/albums/{albumID}/artists/{artistID}/roles/{role}", s.getAlbumArtist())
	r.Get("/albums/{albumID}/disc/{discNumber}/track/{trackNumber}", s.getAlbumTrack())

	r.Get("/album-artists/{albumArtistID}", s.getAlbumArtistByID())

	r.Get("/album-tracks/{albumTrackID}", s.getAlbumTrackByID())

	r.Get("/tracks/{trackID}", s.getTrack())
	r.Put("/tracks/{trackID}", s.updateTrack())
	r.Post("/tracks/{trackID}/ratings", s.addTrackRating())
	r.Get("/tracks/{trackID}/ratings", s.getTrackRatings())

	r.Get("/mediafiles/{mediafileID}", s.getMediaFile())
	r.Get("/mediafiles/{mediafileID}/file", s.getMediaFileFile())

	r.Post("/playlists", s.addPlaylist())
	r.Get("/playlists/{playlistID}", s.getPlaylist())
	r.Put("/playlists/{playlistID}", s.updatePlaylist())
	r.Delete("/playlists/{playlistID}", s.deletePlaylist())
	r.Post("/playlists/{playlistID}/track", s.addPlaylistTrack())

	r.Get("/playlist-tracks/{playlistTrackID}", s.getPlaylistTrack())
	r.Delete("/playlist-tracks/{playlistTrackID}", s.deletePlaylistTrack())

	r.Post("/sync", s.sync())
	r.Post("/sync/play-history", s.syncPlayHistory())

	r.Get("/ratings/{ratingID}", s.getRating())

	r.With(s.authPkg.RoleCheckMiddleware(types.UserRoleAdmin, types.UserRoleImporterWrite)).Post("/import/album", s.importAlbum())

	r.With(s.authPkg.RoleCheckMiddleware(types.UserRoleAdmin)).Get("/admin/entity-events", s.getEntityEvents())

	r.With(s.authPkg.RoleCheckMiddleware(types.UserRoleAdmin, types.UserRoleUsersGet)).Get("/users", s.listUsers())

	return r
}

func (s *Server) apiInfo() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		return Response{
			Payload: ServerApiInfo{
				Name:    s.config.Name,
				Version: config.VERSION,
				ServerInfo: ServerInfo{
					Application: config.SERVER_APP_NAME,
					Status:      HealthzRunning,
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
