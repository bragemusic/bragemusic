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

	resp := types.ListPayload[types.ArtistDetailed]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return 0, err
	}

	return resp.Count, nil
}

func (s ServerClient) CreateArtist(ctx context.Context, artistData types.ArtistBase) error {
	u, err := url.JoinPath(s.baseUrl, "api", "artists")
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, artistData, nil); err != nil {
		return err
	}

	return nil
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

func (s ServerClient) GetArtistTopTracks(ctx context.Context, artistID uuid.UUID) ([]types.TrackDetailed, error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists", artistID.String(), "top-tracks")
	if err != nil {
		return nil, err
	}

	resp := types.ListPayload[types.TrackDetailed]{}
	if err := s.doGetJson(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (s ServerClient) ListArtists(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) (artists []types.ArtistDetailed, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "artists")
	if err != nil {
		return nil, err
	}

	resp := types.ListPayload[types.ArtistDetailed]{}

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
