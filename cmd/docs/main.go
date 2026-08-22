package main

import (
	"fmt"
	"os"

	"github.com/bragemusic/bragemusic/internal/vars"
	"github.com/bragemusic/bragemusic/pkg/config"
	"github.com/bragemusic/bragemusic/pkg/server"
	"github.com/bragemusic/bragemusic/pkg/utils"
	"github.com/swaggest/openapi-go/openapi31"
)

func generateConfigDocs() error {
	serverDocs, err := config.ServerMdDocs()
	if err != nil {
		return err
	}

	clientDocs, err := config.ClientMdDocs()
	if err != nil {
		return err
	}

	header := `
# Configuration

Below is a description of the configuration parameters. You can see the toml and the env ways to enter the data.
`

	out := header + serverDocs + clientDocs

	err = os.WriteFile("docs/config.md", []byte(out), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := generateConfigDocs()
	if err != nil {
		panic(err)
	}

	s := server.Server{}

	refl := openapi31.NewReflector()
	refl.Spec.Info.
		WithTitle("Brage Music API").
		WithVersion(vars.VERSION).
		WithDescription(`
<p>
  Brage Music is a personal music library and streaming platform built around a strongly typed, relational domain model.
</p>

<p>
  The API exposes structured access to albums, artists, tracks, media files, and their linking entities. Rather than embedding relationships directly, Brage Music models connections explicitly using link resources such as <strong>Album Artists</strong> and <strong>Album Tracks</strong>. This allows flexible reuse of tracks across albums, support for multiple artists per album, and consistent metadata composition.
</p>

<h3>Domain Model Principles</h3>

<h4>Separation of Concerns</h4>

<ul>
  <li><strong>Albums</strong> store metadata about a release.</li>
  <li><strong>Artists</strong> store metadata about a performer or contributor.</li>
  <li><strong>Tracks</strong> store metadata about the musical work itself.</li>
  <li><strong>Media Files</strong> represent the physical audio file data.</li>
  <li><strong>Album Artists</strong> link albums to artists.</li>
  <li><strong>Album Tracks</strong> link albums to tracks and define track ordering.</li>
</ul>

<p>
  Relationships are resolved through these link entities rather than nested structures. This enables consistent data modeling and avoids duplication.
</p>

<h4>Multiple Representations</h4>

<p>
  Many resources are available in both:
</p>

<ul>
  <li>Lightweight forms (metadata only)</li>
  <li>Detailed forms (with related artists, tracks, and expanded metadata populated)</li>
</ul>

<p>
  This allows clients to optimize for performance or completeness depending on context.
</p>

<h3>Playback Architecture</h3>

<p>
  Brage Music supports two playback modes:
</p>

<ul>
  <li><strong>Synchronized Mode</strong> – The client maintains a fully synchronized local copy of media files.</li>
  <li><strong>Streaming Mode</strong> – Media files are streamed on demand from the server.</li>
</ul>

<p>
  The API supports both models by exposing media file metadata and streamable audio access where applicable.
</p>

<h3>Intended Usage</h3>

<p>
  The API is designed for trusted clients such as the desktop application and future browser-based builds. It prioritizes predictable REST semantics, explicit relationships, and domain clarity over implicit or loosely structured representations.
</p>
`)

	refl.SpecEns().ComponentsEns().WithSecuritySchemes(map[string]openapi31.SecuritySchemeOrReference{
		"BearerAuth": {
			SecurityScheme: &openapi31.SecurityScheme{
				HTTPBearer: &openapi31.SecuritySchemeHTTPBearer{
					Scheme:       "bearer",
					BearerFormat: utils.Ptr("brg_v1_bNEjO8cs4s8P5rWhH6X4kSMZ2O_g0KnzzW1F4aeyVbw"),
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
		openapi31.Tag{Name: "Admin", Description: utils.Ptr("Server administration endpoints. It includes events, user management, job monitoring and more.")},
		openapi31.Tag{Name: "Albums", Description: utils.Ptr("Everything that has to do with Albums. An Album is an entity that holds metadata about the album itself. It does not know what artist it belongs to, nor the tracks it contains. There are a few functions that returns detailed metadata, then artists and tracks are populated.")},
		openapi31.Tag{Name: "Album Artists", Description: utils.Ptr("Album Artists are the links between an album and an artist. One album can have more than one artist, and each album artist entry is one link between an artist and an album.")},
		openapi31.Tag{Name: "Album Tracks", Description: utils.Ptr("Album Tracks are the links betwen an album and tracks. One track can be on multiple albums, so instead of hardcoding the tracks we have these links. The album track object says album id, track id, and position of the track")},
		openapi31.Tag{Name: "Artists", Description: utils.Ptr("Artists contains all information about an artist. This entity is then linked with albums to create full set of metadata for tracks.")},
		openapi31.Tag{Name: "Import", Description: utils.Ptr("Add new media. Entire albums and single tracks are available. The file will be analysed by the server and added to the database.")},
		openapi31.Tag{Name: "Likes", Description: utils.Ptr("List track likes. Mainly used for syncing.")},
		openapi31.Tag{Name: "Media Files", Description: utils.Ptr("Media Files accounts for the actual music file data. It knows where to find the file and have infromation about what's in it. The media file entries are linked to tracks.")},
		openapi31.Tag{Name: "Playlists", Description: utils.Ptr("Playlists contains the main metadata of a playlist. A Playlist Track entity is used to link an album track to a playlist.")},
		openapi31.Tag{Name: "Playlist Tracks", Description: utils.Ptr("Playlist Tracks are the links between playlists and album tracks. These are the entities that builds the track content of a playlist.")},
		openapi31.Tag{Name: "Ratings", Description: utils.Ptr("Handles user ratings of tracks.")},
		openapi31.Tag{Name: "Sync", Description: utils.Ptr("Endpoints used for syncing local data against the server.")},
		openapi31.Tag{Name: "Tracks", Description: utils.Ptr("Tracks holds the metadata for the actual track. It references one or zero media files. Used together with album tracks and album artists we can build the entire track metadata. That can be retireved by the detailed endpoints.")},
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
