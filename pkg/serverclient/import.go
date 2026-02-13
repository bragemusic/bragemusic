package serverclient

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/bragemusic/core/pkg/types"
)

func (s ServerClient) ImportAlbum(ctx context.Context, r io.Reader, filename string, musicbrainzID *string) error {
	u, err := url.JoinPath(s.baseUrl, "api", "import", "album")
	if err != nil {
		return err
	}

	filename = filepath.Base(filename)

	ia := types.ImportAlbum{
		MusicbrainzID: musicbrainzID,
	}

	return s.importFile(ctx, u, r, filename, ia)
}

func (s ServerClient) importFile(ctx context.Context, u string, r io.Reader, filename string, metadata any) error {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Write multipart body in a goroutine
	go func() {
		defer pw.Close()
		defer writer.Close()

		// ---- metadata part ----
		metaPart, err := writer.CreateFormField("metadata")
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if metadata != nil {
			b, err := json.Marshal(metadata)
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			if _, err := metaPart.Write(b); err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		// ---- file part ----
		filePart, err := writer.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err := io.Copy(filePart, r); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, pr)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
