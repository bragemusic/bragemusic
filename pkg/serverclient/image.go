package serverclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gofrs/uuid/v5"
)

func (s ServerClient) UploadArtistImage(ctx context.Context, artistID uuid.UUID, img FileUpload) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "artists", artistID.String())
	if err != nil {
		return err
	}

	return s.uploadImage(ctx, u, img)
}

func (s ServerClient) UploadAlbumImage(ctx context.Context, id uuid.UUID, img FileUpload) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "albums", id.String())
	if err != nil {
		return err
	}

	return s.uploadImage(ctx, u, img)
}

func (s ServerClient) UploadPlaylistImage(ctx context.Context, id uuid.UUID, img FileUpload) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "playlists", id.String())
	if err != nil {
		return err
	}

	return s.uploadImage(ctx, u, img)
}

func (s ServerClient) UploadUserImage(ctx context.Context, r io.Reader, filename string) error {
	u, err := url.JoinPath(s.baseUrl, "api", "users", "me", "img")
	if err != nil {
		return err
	}

	return s.uploadImageFromReader(ctx, u, filename, r)
}

func (s ServerClient) uploadImageFromReader(ctx context.Context, u string, filename string, r io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, r)
	if err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
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

func (s ServerClient) uploadImage(ctx context.Context, u string, img FileUpload) error {
	if len(img.Data) == 0 {
		return errors.New("empty file")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", img.Name)
	if err != nil {
		return err
	}

	_, err = part.Write(img.Data)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, u, body)
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
