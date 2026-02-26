package serverclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

func (s ServerClient) DownloadMediaFilePart(ctx context.Context, id uuid.UUID, startByte, endByte int, w io.Writer) (finalBytes bool, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "mediafiles", id.String(), "file")
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startByte, endByte))

	resp, err := s.do(ctx, req)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
			return true, nil
		} else {
			return false, err
		}
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return false, err
	}

	return finalBytes, nil
}
