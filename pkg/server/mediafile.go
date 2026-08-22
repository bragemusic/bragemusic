package server

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
)

func (s *Server) mediafileRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/{mediafileID}", s.getMediaFile(), nil, routes.RouteMeta{
			Summary:             "Retrieve media file metadata by ID.",
			Description:         "Returns metadata about the specified media file.",
			ExpectedDescription: "Metadata about the media file",
			Tags:                []string{"Media Files"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{mediafileID}/file", s.getMediaFileFile(), nil, routes.RouteMeta{
			Summary:             "Retrieve a media file by ID.",
			Description:         "Returns the data of the actual music file found behind the given ID",
			ExpectedDescription: "File data",
			Tags:                []string{"Media Files"},
			Errors: []routes.RouteErrorMeta{{
				Description: "Range Not Satisfiable",
				Status:      416,
			}},
			ExpectedStatus: http.StatusOK,
		}),
	}
}

func (s *Server) getMediaFile() routes.RouteFunc[ReqMediaFilesGet, types.MediaFile] {
	return func(ctx context.Context, req ReqMediaFilesGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.MediaFile], err error) {
		mf, err := s.mediamgr.GetMediaFile(ctx, req.MediafileID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.MediaFile]{
			Payload: mf,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) getMediaFileFile() routes.RouteFunc[ReqMediaFilesGetFile, types.NoResponse] {
	return func(ctx context.Context, req ReqMediaFilesGetFile, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		mf, err := s.mediamgr.GetMediaFile(ctx, req.MediafileID)
		if err != nil {
			return resp, err
		}

		path := filepath.Join(s.config.Paths.MusicDir, mf.Filename())

		f, err := os.Open(path)
		if err != nil {
			return resp, err
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return resp, err
		}

		// TODO: detect mime type properly instead of hardcoding
		w.Header().Set("Content-Type", "audio/flac")

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)

		return types.Response[types.NoResponse]{}, nil
	}
}
