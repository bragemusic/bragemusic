package serverclient

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) AddPlaylist(ctx context.Context, playlist types.Playlist) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists")
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, playlist, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) AddSmartPlaylist(ctx context.Context, playlist types.PlaylistBase, filter types.TrackFilter) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists")
	if err != nil {
		return err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return err
	}

	q := ur.Query()
	q.Set("type", "smart")

	ur.RawQuery = q.Encode()

	req := types.ReqPlaylistsAdd{
		PlaylistBase: playlist,
		Filter:       filter,
	}

	if err := s.doPostJson(ctx, ur.String(), req, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", playlistID.String(), "track")
	if err != nil {
		return err
	}

	pt := types.PlaylistTrackReq{
		AlbumID: albumID,
		TrackID: trackID,
	}

	if err := s.doPostJson(ctx, u, pt, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) CountPlaylists(ctx context.Context) (int, error) {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists")
	if err != nil {
		return 0, err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return 0, err
	}

	q := ur.Query()
	q.Set("count", "true")

	ur.RawQuery = q.Encode()

	resp := types.ListPayload[types.Playlist]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (s ServerClient) CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error) {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", playlistID.String(), "tracks")
	if err != nil {
		return 0, err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return 0, err
	}

	q := ur.Query()
	q.Set("count", "true")

	ur.RawQuery = q.Encode()

	resp := types.ListPayload[types.TrackDetailed]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (s ServerClient) DeletePlaylist(ctx context.Context, id uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", id.String())
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlist-tracks", id.String())
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) GetPlaylist(ctx context.Context, id uuid.UUID) (playlist types.Playlist, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", id.String())
	if err != nil {
		return types.Playlist{}, err
	}

	if err := s.doGetJson(ctx, u, &playlist); err != nil {
		return types.Playlist{}, err
	}

	return playlist, nil
}

func (s ServerClient) GetPlaylistTrack(ctx context.Context, id uuid.UUID) (playlistTrack types.PlaylistTrack, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "playlist-tracks", id.String())
	if err != nil {
		return types.PlaylistTrack{}, err
	}

	if err := s.doGetJson(ctx, u, &playlistTrack); err != nil {
		return types.PlaylistTrack{}, err
	}

	return playlistTrack, nil
}

func (s ServerClient) GetSmartPlaylistArtist(ctx context.Context, id uuid.UUID) (artist types.SmartPlaylistArtist, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "smart-playlists", "artist", id.String())
	if err != nil {
		return types.SmartPlaylistArtist{}, err
	}

	if err := s.doGetJson(ctx, u, &artist); err != nil {
		return types.SmartPlaylistArtist{}, err
	}

	return artist, nil
}

func (s ServerClient) GetSmartPlaylistContent(ctx context.Context, id uuid.UUID) (content types.SmartPlaylistContent, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "smart-playlists", "content", id.String())
	if err != nil {
		return types.SmartPlaylistContent{}, err
	}

	if err := s.doGetJson(ctx, u, &content); err != nil {
		return types.SmartPlaylistContent{}, err
	}

	return content, nil
}

func (s ServerClient) DownloadPlaylistImage(ctx context.Context, id uuid.UUID, size imagemagick.ImageSize, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "playlists", id.String(), fmt.Sprintf("%d.jpg", size))
	if err != nil {
		return err
	}

	return s.downloadFile(ctx, u, w)
}

func (s ServerClient) ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error) {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists")
	if err != nil {
		return nil, err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	q := ur.Query()
	q.Set("includePublic", string("true"))
	q.Set("sortBy", string(sortBy))
	q.Set("sortOrder", string(sortOrder))

	ur.RawQuery = q.Encode()
	resp := types.ListPayload[types.Playlist]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (s ServerClient) ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error) {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", playlistID.String(), "tracks")
	if err != nil {
		return nil, err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	q := ur.Query()
	q.Set("sortBy", string(sortBy))
	q.Set("sortOrder", string(sortOrder))

	ur.RawQuery = q.Encode()
	resp := types.ListPayload[types.TrackDetailed]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (s ServerClient) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", id.String())
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, data, nil); err != nil {
		return err
	}

	return nil
}
