package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) importRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("POST", "/album", s.importAlbum(), []types.UserRole{types.UserRoleAdmin, types.UserRoleImporterWrite}, routes.RouteMeta{
			Summary:             "Import an entire album",
			Description:         "Imports an entire album. Must be a zip file.",
			ExpectedDescription: "Import queued",
			Tags:                []string{"Import"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusCreated,
		}),
	}
}

func (s *Server) importAlbum() routes.RouteFunc[ReqImportAlbum, types.NoResponse] {
	return func(ctx context.Context, req ReqImportAlbum, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
		if err != nil {
			return resp, err
		}

		// Get the file from the form input "file"
		file, header, err := r.FormFile("file")
		if err != nil {
			return resp, err
		}
		defer file.Close()

		path := filepath.Join(s.config.Paths.ImportDir, header.Filename)

		// Create the file on the server
		dst, err := os.Create(path)
		if err != nil {
			return resp, err
		}
		defer dst.Close()

		// Copy the uploaded file's content to the destination file
		if _, err = io.Copy(dst, file); err != nil {
			return resp, err
		}
		metaStr := r.FormValue("metadata")
		if metaStr == "" {
			return resp, errors.New("metadata is required")
		}

		var meta types.ImportAlbum
		if err = json.Unmarshal([]byte(metaStr), &meta); err != nil {
			return resp, err
		}

		if err = s.importer.AddImportEntry(ctx, header.Filename, types.ImportTypeAlbum, user.ID, meta.MusicbrainzID); err != nil {
			return resp, err
		}

		if err = s.jobmgr.RunJob(ctx, types.JobImporterRun); err != nil {
			return resp, err
		}

		if err = s.jobmgr.RunJob(ctx, types.JobMetaSyncRun); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusCreated,
		}, nil
	}
}
