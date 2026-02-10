package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/types"
)

func (s Server) importAlbum() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		// ctx := r.Context()

		err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
		if err != nil {
			return Response{}, err
		}

		// Get the file from the form input "file"
		file, header, err := r.FormFile("file")
		if err != nil {
			return Response{}, err
		}
		defer file.Close()

		path := filepath.Join(s.config.Paths.ImportDir, header.Filename)

		// Create the file on the server
		dst, err := os.Create(path)
		if err != nil {
			return Response{}, err
		}
		defer dst.Close()

		// Copy the uploaded file's content to the destination file
		if _, err = io.Copy(dst, file); err != nil {
			return Response{}, err
		}
		metaStr := r.FormValue("metadata")
		if metaStr == "" {
			return Response{}, err
		}

		var meta types.ImportAlbum
		if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
			return Response{}, err
		}

		if meta.MusicbrainzID != nil {
			fmt.Println(*meta.MusicbrainzID)
		}

		// switch imageType {
		// case ArtistImage:
		// 	if err = s.mediamgr.AddArtistImage(ctx, orgImgPath, assetID); err != nil {
		// 		return Response{}, err
		// 	}
		// case AlbumImage:
		// 	if err = s.mediamgr.AddAlbumImage(ctx, orgImgPath, assetID); err != nil {
		// 		return Response{}, err
		// 	}
		// case PlaylistImage:
		// 	if err = s.mediamgr.AddPlaylistImage(ctx, orgImgPath, assetID); err != nil {
		// 		return Response{}, err
		// 	}
		// }

		return Response{Status: http.StatusCreated}, nil
	},
	)
}
