package serverclient

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) DownloadTrackFile(ctx context.Context, trackID string, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "tracks", trackID, "file")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s ServerClient) GetTrack(ctx context.Context, trackID string) (track types.Track, err error) {
	u, err := url.JoinPath(s.baseUrl, "tracks", trackID)
	if err != nil {
		return types.Track{}, err
	}

	if err := s.doGetJson(ctx, u, &track); err != nil {
		return types.Track{}, err
	}

	return track, nil
}

func (s ServerClient) ListTracksByAlbum(ctx context.Context, albumID string) (tracks []types.Track, err error) {
	u, err := url.JoinPath(s.baseUrl, "albums", albumID, "tracks")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}
