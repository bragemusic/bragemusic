package wiki

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const (
	imageUrlFormat = "https://commons.wikimedia.org/wiki/Special:FilePath/%s"
)

type Wiki struct {
	email          string
	preferredLangs []string
	log            *slog.Logger
}

func (w Wiki) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", fmt.Sprintf("brage/1.0.0 (https://example.com/contact; %s)", w.email))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (w Wiki) DownloadFile(ctx context.Context, fromUrl, filename string) error {
	req, err := http.NewRequest("GET", fromUrl, nil)
	if err != nil {
		return err
	}

	resp, err := w.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func New(email string, slogHandler slog.Handler) Wiki {
	return Wiki{
		preferredLangs: []string{"en", "sv"},
		email:          email,
		log:            slog.New(slogHandler).With("service", "wiki"),
	}
}
