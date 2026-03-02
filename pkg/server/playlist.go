package server

import (
	"context"
	"net/http"

	"github.com/bragemusic/core/pkg/routes"
	"github.com/bragemusic/core/pkg/types"
)

func (s *Server) playlistRoutes() []routes.RouteHandler {
	return []routes.RouteHandler{
		routes.New("GET", "/", s.listPlaylists(), nil, routes.RouteMeta{
			Summary:             "List all playlists",
			Description:         "Returns a collection of all playlists available to the client. Each playlist contains metadata only and does not include its tracks.",
			ExpectedDescription: "A list of playlist metadata",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("POST", "/", s.addPlaylist(), nil, routes.RouteMeta{
			Summary:             "Create a new playlist",
			Description:         "Creates a new empty playlist. The playlist is created without any tracks and can be populated through the playlist track endpoints.",
			ExpectedDescription: "Playlist successfully created",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusCreated,
		}),
		routes.New("GET", "/{playlistID}", s.getPlaylist(), nil, routes.RouteMeta{
			Summary:             "Get a playlist by ID",
			Description:         "Returns metadata for the specified playlist. This endpoint does not include the playlist's tracks.",
			ExpectedDescription: "Playlist metadata",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
		routes.New("PUT", "/{playlistID}", s.updatePlaylist(), nil, routes.RouteMeta{
			Summary:             "Update a playlist by ID",
			Description:         "Updates the metadata of the specified playlist. This does not modify the playlist's tracks.",
			ExpectedDescription: "Playlist successfully updated",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("DELETE", "/{playlistID}", s.deletePlaylist(), nil, routes.RouteMeta{
			Summary:             "Delete a playlist by ID",
			Description:         "Deletes the specified playlist and removes all associated playlist track links.",
			ExpectedDescription: "Playlist successfully deleted",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusNoContent,
		}),
		routes.New("POST", "/{playlistID}/track", s.addPlaylistTrack(), nil, routes.RouteMeta{
			Summary:             "Add a track to a playlist",
			Description:         "Creates a playlist track link between the specified playlist and a track. The link defines the track’s position within the playlist.",
			ExpectedDescription: "Playlist track successfully created",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusCreated,
		}),
		routes.New("GET", "/{playlistID}/tracks", s.listPlaylistTracks(), nil, routes.RouteMeta{
			Summary:             "List tracks in a playlist",
			Description:         "Returns all playlist track entries for the specified playlist, including track ordering and associated track metadata.",
			ExpectedDescription: "A list of playlist tracks",
			Tags:                []string{"Playlists"},
			Errors:              []routes.RouteErrorMeta{},
			ExpectedStatus:      http.StatusOK,
		}),
	}
}

func (s *Server) addPlaylist() routes.RouteFunc[ReqPlaylistsAdd, types.NoResponse] {
	return func(ctx context.Context, req ReqPlaylistsAdd, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.AddPlaylist(ctx, req.PlaylistBase, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusCreated,
		}, nil
	}
}

func (s *Server) addPlaylistTrack() routes.RouteFunc[ReqPlaylistsAddTrack, types.NoResponse] {
	return func(ctx context.Context, req ReqPlaylistsAddTrack, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.AddPlaylistTrack(ctx, req.PlaylistID, req.AlbumID, req.TrackID, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusCreated,
		}, nil
	}
}

func (s *Server) deletePlaylist() routes.RouteFunc[ReqPlaylistsGet, types.NoResponse] {
	return func(ctx context.Context, req ReqPlaylistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		err = s.mediamgr.DeletePlaylist(ctx, req.PlaylistID, user.ID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}

func (s *Server) getPlaylist() routes.RouteFunc[ReqPlaylistsGet, types.Playlist] {
	return func(ctx context.Context, req ReqPlaylistsGet, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.Playlist], err error) {
		plist, err := s.mediamgr.GetPlaylist(ctx, req.PlaylistID, user.ID)
		if err != nil {
			return resp, err
		}

		return types.Response[types.Playlist]{
			Payload: plist,
			Status:  http.StatusOK,
		}, nil
	}
}

func (s *Server) listPlaylists() routes.RouteFunc[ReqListPlaylists, types.ListPayload[types.Playlist]] {
	return func(ctx context.Context, req ReqListPlaylists, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.Playlist]], err error) {
		cnt, err := s.mediamgr.CountPlaylists(ctx, user.ID)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.Playlist]]{
				Payload: types.ListPayload[types.Playlist]{
					Items: nil,
					Count: cnt,
				},
				Status: http.StatusOK,
			}, nil
		}

		plists, err := s.mediamgr.ListPlaylists(ctx, user.ID, req.IncludePublic, req.SortBy, req.SortOrder)
		if err != nil {
			return resp, err
		}

		return types.Response[types.ListPayload[types.Playlist]]{
			Payload: types.ListPayload[types.Playlist]{
				Items: plists,
				Count: cnt,
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) listPlaylistTracks() routes.RouteFunc[ReqListPlaylistTracks, types.ListPayload[types.TrackDetailed]] {
	return func(ctx context.Context, req ReqListPlaylistTracks, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.ListPayload[types.TrackDetailed]], err error) {
		cnt, err := s.mediamgr.CountPlaylistTracks(ctx, req.PlaylistID, user.ID)
		if err != nil {
			return resp, err
		}

		if req.Count {
			return types.Response[types.ListPayload[types.TrackDetailed]]{
				Payload: types.ListPayload[types.TrackDetailed]{
					Items: nil,
					Count: cnt,
				},
				Status: http.StatusOK,
			}, nil
		}

		tracks, err := s.mediamgr.ListPlaylistTracks(ctx, req.PlaylistID, user.ID, req.SortBy, req.SortOrder)
		if err != nil {
			return resp, err
		}

		return types.Response[types.ListPayload[types.TrackDetailed]]{
			Payload: types.ListPayload[types.TrackDetailed]{
				Items: tracks,
				Count: cnt,
			},
			Status: http.StatusOK,
		}, nil
	}
}

func (s *Server) updatePlaylist() routes.RouteFunc[ReqPlaylistsUpdate, types.NoResponse] {
	return func(ctx context.Context, req ReqPlaylistsUpdate, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[types.NoResponse], err error) {
		if err := s.mediamgr.UpdatePlaylist(ctx, req.PlaylistID, req.PlaylistBase, user.ID); err != nil {
			return resp, err
		}

		return types.Response[types.NoResponse]{
			Payload: types.NoResponse{},
			Status:  http.StatusNoContent,
		}, nil
	}
}
