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

func (s ServerClient) DownloadArtistImage(ctx context.Context, artistID string, size imagemagick.ImageSize, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "artists", artistID, fmt.Sprintf("%d.jpg", size))
	if err != nil {
		return err
	}

	return s.downloadFile(ctx, u, w)
}

func (s ServerClient) GetArtist(ctx context.Context, artistID uuid.UUID) (artist types.Artist, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists", artistID.String())
	if err != nil {
		return types.Artist{}, err
	}

	if err := s.doGetJson(ctx, u, &artist); err != nil {
		return types.Artist{}, err
	}

	return artist, nil
}

func (s ServerClient) ListArtists(ctx context.Context) (artists []types.Artist, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func (s ServerClient) UpdateArtist(ctx context.Context, artistID string, artistData types.Artist) error {
	u, err := url.JoinPath(s.baseUrl, "api", "artists", artistID)
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, artistData, nil); err != nil {
		return err
	}

	return nil
}
