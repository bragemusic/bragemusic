package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) GetTrack(ctx context.Context, trackID uuid.UUID) (track types.Track, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String())
	if err != nil {
		return types.Track{}, err
	}

	if err := s.doGetJson(ctx, u, &track); err != nil {
		return types.Track{}, err
	}

	return track, nil
}

func (s ServerClient) GetTrackRatings(ctx context.Context, trackID uuid.UUID) (ratings []types.Rating, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String(), "ratings")
	if err != nil {
		return nil, err
	}

	if err := s.doGetJson(ctx, u, &ratings); err != nil {
		return nil, err
	}

	return ratings, nil
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

func (s ServerClient) RateTrack(ctx context.Context, trackID uuid.UUID, value int) error {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String(), "ratings")
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, server.RatingReq{Value: value}, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) UpdateTrack(ctx context.Context, id uuid.UUID, data types.TrackUpdate) error {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", id.String())
	if err != nil {
		return err
	}

	if err := s.doPutJson(ctx, u, data, nil); err != nil {
		return err
	}

	return nil
}
