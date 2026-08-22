package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/routes"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/bragemusic/bragemusic/pkg/utils"
)

func (s *Server) artistRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listArtists(), nil, routes.RouteMeta{
			Summary:             "List all artists.",
			Description:         "Returns metadata about all artists.",
			ExpectedDescription: "Metadata about the artists",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("POST", "/", s.createArtist(), nil, routes.RouteMeta{
			Summary:             "Create an artist.",
			Description:         "Create an artist object on the server.",
			ExpectedDescription: "Artist created",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("GET", "/{artistID}", s.getArtist(), nil, routes.RouteMeta{
			Summary:             "Retrieve an artist by ID.",
			Description:         "Returns metadata about the specified artist.",
			ExpectedDescription: "Metadata about the artist",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{artistID}/albums", s.listAlbumsByArtist(), nil, routes.RouteMeta{
			Summary:             "List an artist's albums.",
			Description:         "Returns all albums the selected artist takes part of.",
			ExpectedDescription: "Artist's albums",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{artistID}/albums/featured", s.listFeaturedAlbumsByArtist(), nil, routes.RouteMeta{
			Summary:             "List an artist's featured albums.",
			Description:         "Returns all albums the selected artist is featured on one or more tracks.",
			ExpectedDescription: "Artist's featured albums",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{artistID}/top-tracks", s.getArtistTopTracks(), nil, routes.RouteMeta{
			Summary:             "Retrieve artist top tracks by ID.",
			Description:         "Returns the 10 most played tracks by the logged in user of the specified artist.",
			ExpectedDescription: "Top tracks",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("PUT", "/{artistID}", s.updateArtist(), nil, routes.RouteMeta{
			Summary:             "Update an artist by ID.",
			Description:         "Updates metadata about the specified artist.",
			ExpectedDescription: "Update succeded",
			Tags:                []string{"Artists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
	}
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

func (s *Server) listFeaturedAlbumsByArtist() routes.RouteFunc[ReqArtistsGet, types.ListPayload[types.AlbumDetailed]] {
	return func(ctx context.Context, req ReqArtistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.AlbumDetailed]], err error) {
		albums, err := s.mediamgr.ListFeaturedAlbumsByArtist(ctx, req.ArtistID, database.SortByDate, database.SortAsc)
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

func (s *Server) createArtist() routes.RouteFunc[ReqArtistsCreate, types.NoResponse] {
	return func(ctx context.Context, req ReqArtistsCreate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err = s.mediamgr.CreateArtist(ctx, req.ArtistBase, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{Status: http.StatusNoContent, Payload: types.NoResponse{}}, nil
	}
}
