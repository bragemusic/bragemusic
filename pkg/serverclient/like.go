package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) GetLike(ctx context.Context, id uuid.UUID) (like types.Like, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "likes", id.String())
	if err != nil {
		return types.Like{}, err
	}

	if err := s.doGetJson(ctx, u, &like); err != nil {
		return types.Like{}, err
	}

	return like, nil
}
