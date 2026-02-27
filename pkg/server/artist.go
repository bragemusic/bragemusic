package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func (s *Server) artistRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listArtists(), routes.RouteMeta{
			Summary:             "List all artists.",
			Description:         "Returns metadata about all artists.",
			ExpectedDescription: "Metadata about the artists",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{artistID}", s.getArtist(), routes.RouteMeta{
			Summary:             "Retrieve an artist by ID.",
			Description:         "Returns metadata about the specified artist.",
			ExpectedDescription: "Metadata about the artist",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{artistID}/albums", s.listAlbumsByArtist(), routes.RouteMeta{
			Summary:             "List an artist's albums.",
			Description:         "Returns all albums the selected artist takes part of.",
			ExpectedDescription: "Artist's albums",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{artistID}/top-tracks", s.getArtistTopTracks(), routes.RouteMeta{
			Summary:             "Retrieve artist top tracks by ID.",
			Description:         "Returns the 10 most played tracks by the logged in user of the specified artist.",
			ExpectedDescription: "Top tracks",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("PUT", "/{artistID}", s.updateArtist(), routes.RouteMeta{
			Summary:             "Update an artist by ID.",
			Description:         "Updates metadata about the specified artist.",
			ExpectedDescription: "Update succeded",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
	}
}

func (s *Server) apiArtists() http.Handler {
	r := chi.NewRouter()

	for _, route := range s.artistRoutes() {
		r.Method(route.Method(), route.Path(), route.Handler(s.log, s.errLog, &s.berr))
	}

	return r
}

func (s *Server) getArtist() routes.RouteFunc[ReqArtistsGet, types.Artist] {
	return func(ctx context.Context, req ReqArtistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.Artist], err error) {
		artist, err := s.mediamgr.GetArtist(ctx, req.ArtistID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.Artist]{Status: http.StatusOK, Payload: artist}, nil
	}
}

func (s *Server) getArtistTopTracks() routes.RouteFunc[ReqArtistsGet, types.ListPayload[types.TrackDetailed]] {
	return func(ctx context.Context, req ReqArtistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.TrackDetailed]], err error) {
		tracks, err := s.mediamgr.ListTracksDetailedByArtist(ctx, req.ArtistID, user.ID, database.SortByPlayCount, database.SortDesc, utils.Ptr(10), false)
		if err != nil {
			return types.Response[types.ListPayload[types.TrackDetailed]]{}, err
		}

		return types.Response[types.ListPayload[types.TrackDetailed]]{
			Status: http.StatusOK,
			Payload: types.ListPayload[types.TrackDetailed]{
				Count: 10,
				Items: tracks,
			},
		}, nil
	}
}

func (s *Server) listAlbumsByArtist() routes.RouteFunc[ReqArtistsGet, types.ListPayload[types.AlbumDetailed]] {
	return func(ctx context.Context, req ReqArtistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.AlbumDetailed]], err error) {
		albums, err := s.mediamgr.ListAlbumsByArtist(ctx, req.ArtistID, database.SortByDate, database.SortAsc)
		if err != nil {
			return resp, err
		}

		return types.Response[types.ListPayload[types.AlbumDetailed]]{
			Status: http.StatusOK,
			Payload: types.ListPayload[types.AlbumDetailed]{
				Count: len(albums),
				Items: albums,
			},
		}, nil
	}
}

func (s *Server) listArtists() routes.RouteFunc[ReqList, types.ListPayload[types.ArtistDetailed]] {
	return func(ctx context.Context, req ReqList, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.ArtistDetailed]], err error) {
		cnt, err := s.mediamgr.CountArtists(ctx)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.ArtistDetailed]]{
				Status: http.StatusOK,
				Payload: types.ListPayload[types.ArtistDetailed]{
					Count: cnt,
				},
			}, nil
		}

		artists, err := s.mediamgr.ListArtists(ctx, database.SortByName, database.SortAsc)
		if err != nil {
			return resp, err
		}

		return types.Response[types.ListPayload[types.ArtistDetailed]]{
			Status: http.StatusOK,
			Payload: types.ListPayload[types.ArtistDetailed]{
				Count: cnt,
				Items: artists,
			},
		}, nil
	}
}

func (s *Server) updateArtist() routes.RouteFunc[ReqArtistsUpdate, types.NoResponse] {
	return func(ctx context.Context, req ReqArtistsUpdate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		artist := types.Artist{ArtistBase: req.ArtistBase}

		if err := s.mediamgr.UpdateArtist(ctx, req.ArtistID, artist, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{Status: http.StatusNoContent, Payload: types.NoResponse{}}, nil
	}
}
