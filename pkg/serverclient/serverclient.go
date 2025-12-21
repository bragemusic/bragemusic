package serverclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/types"
)

type ErrStatus struct {
	Status int
}

func (e ErrStatus) Error() string {
	return fmt.Sprintf("server returned status %d", e.Status)
}

type ServerClient struct {
	log       *slog.Logger
	baseUrl   string
	authToken string
	client    *http.Client
}

func (s *ServerClient) SetAuthToken(token string) {
	s.authToken = token
}

func (s ServerClient) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	s.authToken = "brg_v1_za5pi0duRzyvQXAa6x8YJumVEg0UdAAr5qZGfgwnSJ8"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.authToken))

	resp, err := s.client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return resp, ErrStatus{Status: resp.StatusCode}
	}

	return resp, nil
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

func (s ServerClient) CheckStatus(ctx context.Context) (h server.Status, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "status")
	if err != nil {
		return server.Status{}, err
	}

	if err := s.doGetJson(ctx, u, &h); err != nil {
		return server.Status{}, err
	}

	return h, nil
}

func (s ServerClient) GetUser(ctx context.Context) (user types.UserDetails, err error) {
	u, err := url.JoinPath(s.baseUrl, "api", "user")
	if err != nil {
		return types.UserDetails{}, err
	}

	if err := s.doGetJson(ctx, u, &user); err != nil {
		return types.UserDetails{}, err
	}

	return user, nil
}

func New(baseUrl string, slogHandler slog.Handler) ServerClient {
	return ServerClient{
		log:     slog.New(slogHandler).With("service", "server-client"),
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}
