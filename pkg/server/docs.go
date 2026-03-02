package server

import "github.com/swaggest/openapi-go/openapi31"

func (s *Server) APIDocs(refl *openapi31.Reflector) error {
	basePath := "/api"

	for _, r := range s.albumRoutes() {
		if err := r.Docs(refl, basePath+"/albums"); err != nil {
			return err
		}
	}

	for _, r := range s.albumArtistRoutes() {
		if err := r.Docs(refl, basePath+"/album-artists"); err != nil {
			return err
		}
	}

	for _, r := range s.albumTrackRoutes() {
		if err := r.Docs(refl, basePath+"/album-tracks"); err != nil {
			return err
		}
	}

	for _, r := range s.artistRoutes() {
		if err := r.Docs(refl, basePath+"/artists"); err != nil {
			return err
		}
	}

	for _, r := range s.mediafileRoutes() {
		if err := r.Docs(refl, basePath+"/mediafiles"); err != nil {
			return err
		}
	}

	for _, r := range s.trackRoutes() {
		if err := r.Docs(refl, basePath+"/tracks"); err != nil {
			return err
		}
	}

	for _, r := range s.playlistRoutes() {
		if err := r.Docs(refl, basePath+"/playlists"); err != nil {
			return err
		}
	}

	for _, r := range s.playlistTrackRoutes() {
		if err := r.Docs(refl, basePath+"/playlist-tracks"); err != nil {
			return err
		}
	}

	for _, r := range s.syncRoutes() {
		if err := r.Docs(refl, basePath+"/sync"); err != nil {
			return err
		}
	}

	for _, r := range s.ratingRoutes() {
		if err := r.Docs(refl, basePath+"/ratings"); err != nil {
			return err
		}
	}

	for _, r := range s.searchRoutes() {
		if err := r.Docs(refl, basePath+"/search"); err != nil {
			return err
		}
	}

	return nil
}
