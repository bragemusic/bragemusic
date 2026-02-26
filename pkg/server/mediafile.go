package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid/v5"
)

func (s *Server) getMediaFile() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		mediafileID, err := getParameter[uuid.UUID](ctx, "mediafileID")
		if err != nil {
			return Response{}, err
		}

		track, err := s.mediamgr.GetMediaFile(ctx, mediafileID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: track}, nil
	},
	)
}

func (s *Server) getMediaFileFile() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		mediafileID, err := getParameter[uuid.UUID](ctx, "mediafileID")
		if err != nil {
			return Response{}, err
		}

		mf, err := s.mediamgr.GetMediaFile(ctx, mediafileID)
		if err != nil {
			return Response{}, err
		}

		path := filepath.Join(s.config.Paths.MusicDir, mf.Filename())

		f, err := os.Open(path)
		if err != nil {
			return Response{}, err
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return Response{}, err
		}

		// TODO: detect mime type properly instead of hardcoding
		w.Header().Set("Content-Type", "audio/flac")

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)

		return Response{}, nil
	})
}
