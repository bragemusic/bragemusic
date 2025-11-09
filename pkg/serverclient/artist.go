package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) GetArtist(ctx context.Context, artistID string) (artist types.Artist, err error) {
	u, err := url.JoinPath(s.baseUrl, "artists", artistID)
	if err != nil {
		return types.Artist{}, err
	}

	if err := s.doGetJson(ctx, u, &artist); err != nil {
		return types.Artist{}, err
	}

	return artist, nil
}

func (s ServerClient) ListArtists(ctx context.Context) (artists []types.Artist, err error) {
	u, err := url.JoinPath(s.baseUrl, "artists")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}
