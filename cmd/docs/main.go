package main

import (
	"fmt"

	"github.com/bragemusic/core/internal/config"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/swaggest/openapi-go/openapi31"
)

func main() {
	s := server.Server{}

	refl := openapi31.NewReflector()
	refl.Spec.Info.
		WithTitle("Brage Music API").
		WithVersion(">=" + config.VERSION).
		WithDescription("Put something here")

	refl.SpecEns().ComponentsEns().WithSecuritySchemes(map[string]openapi31.SecuritySchemeOrReference{
		"BearerAuth": {
			SecurityScheme: &openapi31.SecurityScheme{
				HTTPBearer: &openapi31.SecuritySchemeHTTPBearer{
					Scheme: "bearer",
				},
			},
		},
	})
	refl.SpecEns().WithSecurity(
		map[string][]string{
			"BearerAuth": {},
		},
	)
	refl.SpecEns().WithTags(
		openapi31.Tag{Name: "Artists", Description: utils.Ptr("Artist management")},
		openapi31.Tag{Name: "Albums", Description: utils.Ptr("Everything that has to do with Albums. An Album is an entity that holds metadata about the album itself. It does not know what artist it belongs to, nor the tracks it contains. There are a few functions that returns detailed metadata, then artists and tracks are populated.")},
		openapi31.Tag{Name: "Album Artists", Description: utils.Ptr("Album Artists are the links between an album and an artist. One album can have more than one artist, and each album artist entry is one link between an artist and an album.")},
		openapi31.Tag{Name: "Album Tracks", Description: utils.Ptr("Album Tracks are the links betwen an album and tracks. One track can be on multiple albums, so instead of hardcoding the tracks we have these links. The album track object says album id, track id, and position of the track")},
	)

	if err := s.APIDocs(refl); err != nil {
		panic(err)
	}
	// fmt.Println(ro.Docs(refl, "/api"))
	schema, err := refl.Spec.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(schema))
}
