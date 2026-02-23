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

func (s ServerClient) CountAlbums(ctx context.Context) (cnt int, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums")
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

	resp := server.ListPayload[types.AlbumDetailed]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (s ServerClient) DownloadAlbumCover(ctx context.Context, albumID uuid.UUID, size imagemagick.ImageSize, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "albums", albumID.String(), fmt.Sprintf("%d.jpg", size))
	if err != nil {
		return err
	}

	return s.downloadFile(ctx, u, w)
}

func (s ServerClient) GetAlbum(ctx context.Context, albumID uuid.UUID) (album types.Album, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID.String())
	if err != nil {
		return types.Album{}, err
	}

	if err := s.doGetJson(ctx, u, &album); err != nil {
		return types.Album{}, err
	}

	return album, nil
}

func (s ServerClient) GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (album types.AlbumDetailed, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID.String(), "detailed")
	if err != nil {
		return types.AlbumDetailed{}, err
	}

	if err := s.doGetJson(ctx, u, &album); err != nil {
		return types.AlbumDetailed{}, err
	}

	return album, nil
}

func (s ServerClient) ListAlbumsByArtist(ctx context.Context, artistID string) (albums []types.Album, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists", artistID, "albums")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &albums); err != nil {
		return nil, err
	}

	return albums, nil
}

func (s ServerClient) GetAlbumArtist(ctx context.Context, albumID, artistID uuid.UUID, role types.ArtistRole) (albumArtist types.AlbumArtist, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID.String(), "artists", artistID.String(), "roles", string(role))
	if err != nil {
		return types.AlbumArtist{}, err
	}

	if err := s.doGetJson(ctx, u, &albumArtist); err != nil {
		return types.AlbumArtist{}, err
	}

	return albumArtist, nil
}

func (s ServerClient) GetAlbumArtistByID(ctx context.Context, id uuid.UUID) (albumArtist types.AlbumArtist, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "album-artists", id.String())
	if err != nil {
		return types.AlbumArtist{}, err
	}

	if err := s.doGetJson(ctx, u, &albumArtist); err != nil {
		return types.AlbumArtist{}, err
	}

	return albumArtist, nil
}

func (s ServerClient) GetAlbumTrack(ctx context.Context, albumID uuid.UUID, discNumber, trackNumber int) (albumTrack types.AlbumTrack, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID.String(), "disc", fmt.Sprint(discNumber), "track", fmt.Sprint(trackNumber))
	if err != nil {
		return types.AlbumTrack{}, err
	}

	if err := s.doGetJson(ctx, u, &albumTrack); err != nil {
		return types.AlbumTrack{}, err
	}

	return albumTrack, nil
}

func (s ServerClient) GetAlbumTrackByID(ctx context.Context, id uuid.UUID) (albumTrack types.AlbumTrack, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "album-tracks", id.String())
	if err != nil {
		return types.AlbumTrack{}, err
	}

	if err := s.doGetJson(ctx, u, &albumTrack); err != nil {
		return types.AlbumTrack{}, err
	}

	return albumTrack, nil
}

func (s ServerClient) UpdateAlbum(ctx context.Context, albumID uuid.UUID, albumData types.AlbumUpdate) error {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID.String())
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, albumData, nil); err != nil {
		return err
	}

	return nil
}
