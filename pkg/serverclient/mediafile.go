package serverclient

import (
	"context"
	"io"
	"net/url"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) GetMediaFile(ctx context.Context, id uuid.UUID) (mediaFile types.MediaFile, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "mediafiles", id.String())
	if err != nil {
		return types.MediaFile{}, err
	}

	if err := s.doGetJson(ctx, u, &mediaFile); err != nil {
		return types.MediaFile{}, err
	}

	return mediaFile, nil
}

func (s ServerClient) DownloadMediaFile(ctx context.Context, id uuid.UUID, w io.Writer) error {
	u, err := url.JoinPath(s.baseUrl, "api", "mediafiles", id.String(), "file")
	if err != nil {
		return err
	}

	return s.downloadFile(ctx, u, w)
}
