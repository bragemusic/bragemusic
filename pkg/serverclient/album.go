package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) GetAlbum(ctx context.Context, albumID string) (album types.Album, err error) {
	u, err := url.JoinPath(s.baseUrl, "albums", albumID)
	if err != nil {
		return types.Album{}, err
	}

	if err := s.doGetJson(ctx, u, &album); err != nil {
		return types.Album{}, err
	}

	return album, nil
}

func (s ServerClient) ListAlbumsByArtist(ctx context.Context, artistID string) (albums []types.Album, err error) {
	u, err := url.JoinPath(s.baseUrl, "artists", artistID, "albums")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &albums); err != nil {
		return nil, err
	}

	return albums, nil
}
