package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) GetTrack(ctx context.Context, trackID string) (track types.Track, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID)
	if err != nil {
		return types.Track{}, err
	}

	if err := s.doGetJson(ctx, u, &track); err != nil {
		return types.Track{}, err
	}

	return track, nil
}

func (s ServerClient) ListTracksByAlbum(ctx context.Context, albumID string) (tracks []types.Track, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID, "tracks")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}
