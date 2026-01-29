package serverclient

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/server"
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

func (s ServerClient) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "playlists", playlistID.String(), "track")
	if err != nil {
		return err
	}

	pt := server.PlaylistTrackReq{
		AlbumID: albumID,
		TrackID: trackID,
	}

	if err := s.doPostJson(ctx, u, pt, nil); err != nil {
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

func (s ServerClient) DownloadPlaylistImage(ctx context.Context, id uuid.UUID, size imagemagick.ImageSize, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "playlists", id.String(), fmt.Sprintf("%d.jpg", size))
	if err != nil {
		return err
	}

	return s.downloadFile(ctx, u, w)
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
