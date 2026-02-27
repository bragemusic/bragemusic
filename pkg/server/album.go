package server

import (
	"encoding/json"
	"net/http"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s *Server) getAlbum() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		album, err := s.mediamgr.GetAlbum(ctx, albumID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: album}, nil
	})
}

func (s *Server) getAlbumDetailed() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		album, err := s.mediamgr.GetAlbumDetailed(ctx, albumID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: album}, nil
	})
}

func (s *Server) listAlbums() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		cnt, err := s.mediamgr.CountAlbums(ctx)
		if err != nil {
			return Response{}, err
		}

		if r.URL.Query().Get("count") == "true" {
			return Response{Status: http.StatusOK, Payload: types.ListPayload[types.AlbumDetailed]{Count: cnt}}, nil
		}

		albums, err := s.mediamgr.ListAlbums(ctx, database.SortByName, database.SortAsc)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: types.ListPayload[types.AlbumDetailed]{Count: cnt, Items: albums}}, nil
	})
}

func (s *Server) getAlbumArtist() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		artistID, err := getParameter[uuid.UUID](ctx, "artistID")
		if err != nil {
			return Response{}, err
		}

		role, err := getParameter[string](ctx, "role")
		if err != nil {
			return Response{}, err
		}

		albumArtist, err := s.mediamgr.GetAlbumArtist(ctx, albumID, artistID, role)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: albumArtist}, nil
	})
}

func (s *Server) getAlbumArtistByID() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumArtistID, err := getParameter[uuid.UUID](ctx, "albumArtistID")
		if err != nil {
			return Response{}, err
		}

		albumArtist, err := s.mediamgr.GetAlbumArtistByID(ctx, albumArtistID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: albumArtist}, nil
	})
}

func (s *Server) getAlbumTrack() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		discNumber, err := getParameter[int](ctx, "discNumber")
		if err != nil {
			return Response{}, err
		}

		trackNumber, err := getParameter[int](ctx, "trackNumber")
		if err != nil {
			return Response{}, err
		}

		albumArtist, err := s.mediamgr.GetAlbumTrack(ctx, albumID, discNumber, trackNumber)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: albumArtist}, nil
	})
}

func (s *Server) getAlbumTrackByID() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		albumTrackID, err := getParameter[uuid.UUID](ctx, "albumTrackID")
		if err != nil {
			return Response{}, err
		}

		albumTrack, err := s.mediamgr.GetAlbumTrackByID(ctx, albumTrackID)
		if err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusOK, Payload: albumTrack}, nil
	})
}

func (s *Server) updateAlbum() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		id, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			return Response{}, err
		}

		album := types.AlbumUpdate{}
		if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.UpdateAlbum(ctx, id, album, user.ID); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
