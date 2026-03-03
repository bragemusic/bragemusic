package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) albumRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listAlbums(), nil, routes.RouteMeta{
			Summary:             "List all albums.",
			Description:         "Returns metadata about all albums.",
			ExpectedDescription: "Metadata about the albums",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{albumID}", s.getAlbum(), nil, routes.RouteMeta{
			Summary:             "Retrieve an album by ID.",
			Description:         "Returns metadata about the specified album.",
			ExpectedDescription: "Metadata about the album",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("PUT", "/{albumID}", s.updateAlbum(), nil, routes.RouteMeta{
			Summary:             "Update an album by ID.",
			Description:         "Updates metadata about the specified album.",
			ExpectedDescription: "Update succeded",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("GET", "/{albumID}/detailed", s.getAlbumDetailed(), nil, routes.RouteMeta{
			Summary:             "Retrieve a detailed album by ID.",
			Description:         "Returns detailed metadata about the specified album.",
			ExpectedDescription: "Metadata about the album",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{albumID}/tracks", s.listAlbumTracks(), nil, routes.RouteMeta{
			Summary:             "List tracks of an album.",
			Description:         "Returns metadata about all the tracks on an album.",
			ExpectedDescription: "Metadata about the tracks",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{albumID}/tracks/{trackID}", s.getTrackDetailed(), nil, routes.RouteMeta{
			Summary:             "Retrieve a track from the album by ID.",
			Description:         "Returns detailed metadata about the wanted track of the specified album.",
			ExpectedDescription: "Metadata about the track",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{albumID}/tracks-detailed", s.listAlbumTracksDetailed(), nil, routes.RouteMeta{
			Summary:             "List detailed tracks of an album.",
			Description:         "Returns detailed metadata about all the tracks on an album.",
			ExpectedDescription: "Metadata about the tracks",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("GET", "/{albumID}/artists/{artistID}/roles/{role}", s.getAlbumArtist(), nil, routes.RouteMeta{
			Summary:             "Retrieve an album artist.",
			Description:         "Returns metadata about the album artist of the specified album.",
			ExpectedDescription: "Metadata about the album artist",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
			Deprecated:          true,
		}),
		routes.New("GET", "/{albumID}/disc/{discNumber}/track/{trackNumber}", s.getAlbumTrack(), nil, routes.RouteMeta{
			Summary:             "Retrieve an album track.",
			Description:         "Returns metadata about the album track of the specified album.",
			ExpectedDescription: "Metadata about the album track",
			Tags:                []string{"Albums"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
			Deprecated:          true,
		}),
	}
}

func (s *Server) getAlbum() routes.RouteFunc[ReqAlbumsGet, types.Album] {
	return func(ctx context.Context, req ReqAlbumsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.Album], err error) {
		album, err := s.mediamgr.GetAlbum(ctx, req.AlbumID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.Album]{Status: http.StatusOK, Payload: album}, nil
	}
}

func (s *Server) getAlbumArtist() routes.RouteFunc[ReqAlbumsAlbumArtistGet, types.AlbumArtist] {
	return func(ctx context.Context, req ReqAlbumsAlbumArtistGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.AlbumArtist], err error) {
		albumArtist, err := s.mediamgr.GetAlbumArtist(ctx, req.AlbumID, req.ArtistID, req.Role)
		if err != nil {
			return resp, err
		}

		return types.Response[types.AlbumArtist]{
			Payload: albumArtist,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) getAlbumDetailed() routes.RouteFunc[ReqAlbumsGet, types.AlbumDetailed] {
	return func(ctx context.Context, req ReqAlbumsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.AlbumDetailed], err error) {
		album, err := s.mediamgr.GetAlbumDetailed(ctx, req.AlbumID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.AlbumDetailed]{
			Payload: album,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) getAlbumTrack() routes.RouteFunc[ReqAlbumsAlbumTrackGet, types.AlbumTrack] {
	return func(ctx context.Context, req ReqAlbumsAlbumTrackGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.AlbumTrack], err error) {
		albumTrack, err := s.mediamgr.GetAlbumTrack(ctx, req.AlbumID, req.DiscNumber, req.TrackNumber)
		if err != nil {
			return resp, err
		}

		return types.Response[types.AlbumTrack]{
			Payload: albumTrack,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) getTrackDetailed() routes.RouteFunc[ReqAlbumsTrackGet, types.TrackDetailed] {
	return func(ctx context.Context, req ReqAlbumsTrackGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.TrackDetailed], err error) {
		track, err := s.mediamgr.GetTrackDetailed(ctx, req.TrackID, req.AlbumID, user.ID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.TrackDetailed]{
			Payload: track,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) listAlbums() routes.RouteFunc[ReqList, types.ListPayload[types.AlbumDetailed]] {
	return func(ctx context.Context, req ReqList, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.AlbumDetailed]], err error) {
		cnt, err := s.mediamgr.CountAlbums(ctx)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.AlbumDetailed]]{
				Payload: types.ListPayload[types.AlbumDetailed]{
					Items: nil,
					Count: cnt,
				},
				Status: http.StatusOK,
			}, nil
		}

		albums, err := s.mediamgr.ListAlbums(ctx, database.SortByName, database.SortAsc)
		if err != nil {
			return resp, err
		}

		return types.Response[types.ListPayload[types.AlbumDetailed]]{
			Payload: types.ListPayload[types.AlbumDetailed]{
				Items: albums,
				Count: cnt,
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) listAlbumTracks() routes.RouteFunc[ReqListTracksOfAlbum, types.ListPayload[types.Track]] {
	return func(ctx context.Context, req ReqListTracksOfAlbum, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.Track]], err error) {
		tracks, err := s.mediamgr.ListTracksByAlbum(ctx, req.AlbumID)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.Track]]{
				Payload: types.ListPayload[types.Track]{
					Items: nil,
					Count: len(tracks),
				},
				Status: http.StatusOK,
			}, nil
		}

		return types.Response[types.ListPayload[types.Track]]{
			Payload: types.ListPayload[types.Track]{
				Items: tracks,
				Count: len(tracks),
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) listAlbumTracksDetailed() routes.RouteFunc[ReqListTracksOfAlbum, types.ListPayload[types.TrackDetailed]] {
	return func(ctx context.Context, req ReqListTracksOfAlbum, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.TrackDetailed]], err error) {
		tracks, err := s.mediamgr.ListTracksDetailedByAlbum(ctx, req.AlbumID, user.ID)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.TrackDetailed]]{
				Payload: types.ListPayload[types.TrackDetailed]{
					Items: nil,
					Count: len(tracks),
				},
				Status: http.StatusOK,
			}, nil
		}

		return types.Response[types.ListPayload[types.TrackDetailed]]{
			Payload: types.ListPayload[types.TrackDetailed]{
				Items: tracks,
				Count: len(tracks),
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) updateAlbum() routes.RouteFunc[ReqAlbumsUpdate, types.NoResponse] {
	return func(ctx context.Context, req ReqAlbumsUpdate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.UpdateAlbum(ctx, req.AlbumID, req.AlbumUpdate, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}
