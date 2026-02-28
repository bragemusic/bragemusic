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
		openapi31.Tag{Name: "Albums", Description: utils.Ptr("Album management")},
		// openapi3.Tag{Name: "Tracks", Description: "Track management"},
		// openapi3.Tag{Name: "Auth", Description: "Authentication endpoints"},
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
