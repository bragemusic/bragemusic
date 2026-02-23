package serverclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) CountArtists(ctx context.Context) (cnt int, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists")
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

	resp := server.Response{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return 0, err
	}

	pl, ok := resp.Payload.(server.ListPayload[types.ArtistDetailed])
	if !ok {
		return 0, errors.New("wrong payload response type")
	}

	return pl.Count, nil
}

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

func (s ServerClient) ListArtists(ctx context.Context) (artists []types.ArtistDetailed, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists")
	if err != nil {
		return nil, err
	}

	resp := server.ListPayload[types.ArtistDetailed]{}

	if err := s.doGetJson(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (s ServerClient) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	u, err := url.JoinPath(s.baseUrl, "api", "artists", artistID.String())
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, artistData, nil); err != nil {
		return err
	}

	return nil
}
