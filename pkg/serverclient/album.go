package serverclient

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) DownloadAlbumCover(ctx context.Context, albumID string, size imagemagick.ImageSize, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "albums", albumID, fmt.Sprintf("%d.jpg", size))
	if err != nil {
		return err
	}

	return s.downloadFile(ctx, u, w)
}

func (s ServerClient) GetAlbum(ctx context.Context, albumID string) (album types.Album, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID)
	if err != nil {
		return types.Album{}, err
	}

	if err := s.doGetJson(ctx, u, &album); err != nil {
		return types.Album{}, err
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

func (s ServerClient) UpdateAlbum(ctx context.Context, albumID string, albumData types.AlbumUpdate) error {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID)
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, albumData, nil); err != nil {
		return err
	}

	return nil
}
