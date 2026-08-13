package analyser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/internalusers"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (a Analyser) CheckAvailability(ctx context.Context) error {
	u, err := url.JoinPath(a.analyserBaseURL, "healthz")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	_, err = a.do(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (a Analyser) RunTrackAnalysisJob(ctx context.Context) error {
	a.log.InfoContext(ctx, "starting track analysis check")

	for {
		trackID, found, err := a.db.GetUnanalysedTrack(ctx)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		if err = a.CheckAvailability(ctx); err != nil {
			return errors.New("track analyser is not available")
		}

		if err = a.RunTrackAnalysis(ctx, trackID); err != nil {
			a.log.ErrorContext(ctx, "could not analyse track", "error", err.Error(), "track_id", trackID)
			continue
		}
	}

	a.log.InfoContext(ctx, "track analysis check done")
	return nil
}

func (a Analyser) RunTrackAnalysis(ctx context.Context, trackID uuid.UUID) error {
	track, err := a.db.GetTrackFromID(ctx, trackID)
	if err != nil {
		return a.berr.DatabaseError(err, types.EntityTrack, &trackID)
	}

	if track.MediaFile == nil {
		return a.berr.NoMediaFile(errors.New("could not start analysis"), track.Title)
	}

	mf, err := a.db.GetMediaFile(ctx, *track.MediaFile)
	if err != nil {
		return a.berr.DatabaseError(err, types.EntityMediaFile, track.MediaFile)
	}

	if mf.Codec != types.CodecFlac {
		return errors.New("only flac supported in analysis")
	}

	a.log.InfoContext(ctx, "analysing track", "track_name", track.Title, "track_id", track.ID, "mediafile_id", mf.ID)

	filename := filepath.Join(a.musicDirPath, fmt.Sprintf("%s.%s", mf.ID.String(), mf.Codec))

	res, err := a.runMediafileAnalysis(ctx, filename)
	if err != nil {
		return err
	}

	a.log.InfoContext(ctx, "analysis finished", "track_name", track.Title, "track_id", track.ID, "mediafile_id", mf.ID)

	trackAnalysis := types.TrackAnalysis{
		ID:                   trackID,
		TrackAnalysisResults: res,
	}

	if err = a.db.AddTrackAnalysis(ctx, trackAnalysis, internalusers.TrackAnalyser); err != nil {
		return a.berr.DatabaseError(err, types.EntityTrackAnalysis, &track.ID)
	}

	return nil
}

func (a Analyser) runMediafileAnalysis(ctx context.Context, filename string) (types.TrackAnalysisResults, error) {
	r, err := os.Open(filename)
	if err != nil {
		return types.TrackAnalysisResults{}, err
	}
	defer r.Close()

	u, err := url.JoinPath(a.analyserBaseURL, "run")
	if err != nil {
		return types.TrackAnalysisResults{}, err
	}

	results, err := a.uploadFile(ctx, filename, r, u)
	if err != nil {
		return types.TrackAnalysisResults{}, err
	}

	return results, nil
}

func (a Analyser) uploadFile(ctx context.Context, filename string, r io.Reader, u string) (types.TrackAnalysisResults, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Write multipart body in a goroutine
	go func() {
		defer pw.Close()
		defer writer.Close()
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
		return types.TrackAnalysisResults{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := a.do(ctx, req)
	if err != nil {
		return types.TrackAnalysisResults{}, err
	}
	defer resp.Body.Close()

	var results types.TrackAnalysisResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return types.TrackAnalysisResults{}, err
	}

	return results, nil
}
