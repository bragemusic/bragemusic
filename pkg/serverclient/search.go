package serverclient

import (
	"context"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) SearchFull(ctx context.Context, searchTerm string) (si []types.SearchItem, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "search")
	if err != nil {
		return nil, err
	}

	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	q := ur.Query()
	q.Set("q", searchTerm)

	ur.RawQuery = q.Encode()

	resp := types.ListPayload[types.SearchItem]{}

	if err := s.doGetJson(ctx, ur.String(), &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}
