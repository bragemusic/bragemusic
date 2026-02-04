package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (s Server) getAlbum() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		album, err := s.mediamgr.GetAlbum(ctx, albumID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, album, nil
	})
}

func (s Server) listAlbums() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		artistID, err := getParameter[uuid.UUID](ctx, "artistID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		albums, err := s.mediamgr.ListAlbumsByArtist(ctx, artistID, database.SortByDate, database.SortAsc)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "artistID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, albums, nil
	})
}

func (s Server) getAlbumArtist() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumID := chi.URLParamFromCtx(ctx, "albumID")
		if albumID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "albumID",
				err:   errors.New("could not parse albumID"),
			}
		}

		artistID := chi.URLParamFromCtx(ctx, "artistID")
		if albumID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "artistID",
				err:   errors.New("could not parse artistID"),
			}
		}

		role := chi.URLParamFromCtx(ctx, "role")
		if albumID == "" {
			return http.StatusBadRequest, nil, ErrIDNotFound{
				idKey: "role",
				err:   errors.New("could not parse role"),
			}
		}

		albumArtist, err := s.mediamgr.GetAlbumArtist(ctx, albumID, artistID, role)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, albumArtist, nil
	})
}

func (s Server) getAlbumArtistByID() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumArtistID, err := getParameter[uuid.UUID](ctx, "albumArtistID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		albumArtist, err := s.mediamgr.GetAlbumArtistByID(ctx, albumArtistID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumArtistID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, albumArtist, nil
	})
}

func (s Server) getAlbumTrack() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumID, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		discNumber, err := getParameter[int](ctx, "discNumber")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		trackNumber, err := getParameter[int](ctx, "trackNumber")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		albumArtist, err := s.mediamgr.GetAlbumTrack(ctx, albumID, discNumber, trackNumber)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, albumArtist, nil
	})
}

func (s Server) getAlbumTrackByID() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		albumTrackID, err := getParameter[uuid.UUID](ctx, "albumTrackID")
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		albumTrack, err := s.mediamgr.GetAlbumTrackByID(ctx, albumTrackID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return http.StatusBadRequest, nil, ErrIDNotFound{
					idKey: "albumArtistID",
					err:   err,
				}
			} else {
				return http.StatusInternalServerError, nil, err
			}
		}

		return http.StatusOK, albumTrack, nil
	})
}

func (s Server) updateAlbum() http.HandlerFunc {
	return s.handle(func(w http.ResponseWriter, r *http.Request) (Response, error) {
		ctx := r.Context()

		id, err := getParameter[uuid.UUID](ctx, "albumID")
		if err != nil {
			return Response{}, err
		}

		album := types.AlbumUpdate{}
		if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
			return Response{}, err
		}

		if err := s.mediamgr.UpdateAlbum(ctx, id, album); err != nil {
			return Response{}, err
		}

		return Response{Status: http.StatusNoContent}, nil
	})
}
