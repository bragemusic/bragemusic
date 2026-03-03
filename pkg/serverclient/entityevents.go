package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) ListEntityEvents(ctx context.Context) (events []types.EntityEvent, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "admin", "entity-events")
	if err != nil {
		return nil, err
	}

	resp := types.ListPayload[types.EntityEvent]{}

	if err := s.doGetJson(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}
