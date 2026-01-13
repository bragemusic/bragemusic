package serverclient

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (s ServerClient) UploadArtistImage(ctx context.Context, artistID string, img ImageUpload) error {
	u, err := url.JoinPath(s.baseUrl, "api", "img", "artists", artistID)
	if err != nil {
		return err
	}

	return s.uploadImage(ctx, u, img)
}

func (s ServerClient) uploadImage(ctx context.Context, u string, img ImageUpload) error {
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
