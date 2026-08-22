package analyser

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/bragerr"
	"github.com/bragemusic/bragemusic/pkg/database"
)

type Analyser struct {
	analyserBaseURL string
	musicDirPath    string

	db     database.DatabaseFace
	client *http.Client
	berr   bragerr.BragErrFactory
	log    *slog.Logger
}

func (a Analyser) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := a.client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return resp, errors.New(resp.Status)
	}

	return resp, nil
}

func New(baseURL string, musicDirPath string, db database.DatabaseFace, slogHandler slog.Handler) Analyser {
	return Analyser{
		analyserBaseURL: baseURL,
		musicDirPath:    musicDirPath,
		db:              db,
		client:          http.DefaultClient,
		berr:            bragerr.NewFactory("analyser"),
		log:             slog.New(slogHandler).With("service", "analyser"),
	}
}
