package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) AddPlayCount(ctx context.Context, trackID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String(), "play-history")
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, nil, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) CountTracks(ctx context.Context) (cnt int, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks")
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

	resp := types.ListPayload[types.TrackDetailed]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return 0, err
	}

	return resp.Count, nil
}

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

func (s ServerClient) GetTrackDetailed(ctx context.Context, trackID, albumID uuid.UUID) (track types.TrackDetailed, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "albums", albumID.String(), "tracks", trackID.String())
	if err != nil {
		return types.TrackDetailed{}, err
	}

	if err := s.doGetJson(ctx, u, &track); err != nil {
		return types.TrackDetailed{}, err
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

func (s ServerClient) RateTrack(ctx context.Context, trackID uuid.UUID, value int) error {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String(), "ratings")
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, types.RatingReq{Value: value}, nil); err != nil {
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

func (s ServerClient) LikeTrack(ctx context.Context, trackID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String(), "like")
	if err != nil {
		return err
	}

	if err := s.doPostJson(ctx, u, nil, nil); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) UnlikeTrack(ctx context.Context, trackID uuid.UUID) error {
	u, err := url.JoinPath(s.baseUrl, "api", "tracks", trackID.String(), "like")
	if err != nil {
		return err
	}

	if err := s.doDelete(ctx, u); err != nil {
		return err
	}

	return nil
}
