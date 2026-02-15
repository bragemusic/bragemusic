package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type ImageType string

const (
	ArtistImage   ImageType = "artist"
	AlbumImage    ImageType = "album"
	PlaylistImage ImageType = "playlist"
)

func (s *Server) getImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "*")
		fp := strings.TrimPrefix(r.URL.Path, pathPrefix)

		filename := filepath.Join(s.config.Paths.ImageDir, fp)

		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, filename)
	}
}

func (s *Server) addImage(imageType ImageType) http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		var assetID uuid.UUID
		var err error

		switch imageType {
		case ArtistImage:
			assetID, err = getParameter[uuid.UUID](ctx, "artistID")
			if err != nil {
				return Response{}, err
			}
		case AlbumImage:
			assetID, err = getParameter[uuid.UUID](ctx, "albumID")
			if err != nil {
				return Response{}, err
			}
		case PlaylistImage:
			assetID, err = getParameter[uuid.UUID](ctx, "playlistID")
			if err != nil {
				return Response{}, err
			}
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		err = r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
		if err != nil {
			return Response{}, err
		}

		// Get the file from the form input "file"
		file, _, err := r.FormFile("file")
		if err != nil {
			return Response{}, err
		}
		defer file.Close()

		tempFolder, err := os.MkdirTemp(os.TempDir(), "brage-img")
		if err != nil {
			return Response{}, err
		}
		defer os.RemoveAll(tempFolder)

		orgImgPath := filepath.Join(tempFolder, fmt.Sprintf("%s.jpg", assetID.String()))

		// Create the file on the server
		dst, err := os.Create(orgImgPath)
		if err != nil {
			return Response{}, err
		}
		defer dst.Close()

		// Copy the uploaded file's content to the destination file
		if _, err = io.Copy(dst, file); err != nil {
			return Response{}, err
		}

		switch imageType {
		case ArtistImage:
			if err = s.mediamgr.AddArtistImage(ctx, orgImgPath, assetID, user.ID); err != nil {
				return Response{}, err
			}
		case AlbumImage:
			if err = s.mediamgr.AddAlbumImage(ctx, orgImgPath, assetID); err != nil {
				return Response{}, err
			}
		case PlaylistImage:
			if err = s.mediamgr.AddPlaylistImage(ctx, orgImgPath, assetID); err != nil {
				return Response{}, err
			}
		}

		return Response{Status: http.StatusCreated}, nil
	},
	)
}
