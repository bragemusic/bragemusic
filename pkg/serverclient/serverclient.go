package serverclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type ServerClient struct {
	log     *slog.Logger
	baseUrl string
	client  *http.Client
}

func (s ServerClient) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

func (s ServerClient) doGetJson(ctx context.Context, u string, target any) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) doPostJson(ctx context.Context, u string, payload any, target any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}

func (s ServerClient) downloadFile(ctx context.Context, u string, w io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := s.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func New(baseUrl string, slogHandler slog.Handler) ServerClient {
	return ServerClient{
		log:     slog.New(slogHandler).With("service", "server-client"),
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}
