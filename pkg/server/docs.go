package server

import "github.com/swaggest/openapi-go/openapi31"

func (s *Server) APIDocs(refl *openapi31.Reflector) error {
	basePath := "/api"

	for _, r := range s.artistRoutes() {
		if err := r.Docs(refl, basePath+"/artists"); err != nil {
			return err
		}
	}

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

	return nil
}
