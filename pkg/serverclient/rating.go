package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) GetRating(ctx context.Context, id uuid.UUID) (rating types.Rating, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "ratings", id.String())
	if err != nil {
		return types.Rating{}, err
	}

	if err := s.doGetJson(ctx, u, &rating); err != nil {
		return types.Rating{}, err
	}

	return rating, nil
}
