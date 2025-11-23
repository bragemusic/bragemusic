package acoustid

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os/exec"
	"time"
)

const (
	apiUrl        = "https://api.acoustid.org/v2/lookup"
	scoreLimit    = 0.96
	minReqTimeout = 1100 * time.Millisecond
)

type fpCalcResp struct {
	Duration    float32 `json:"duration"`
	Fingerprint string  `json:"fingerprint"`
}

type Response struct {
	Results []Results `json:"results"`
}

type Results struct {
	Releases []Release `json:"releases"`
	Score    float32   `json:"score"`
}

type Release struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
	Mediums []Medium `json:"mediums"`
	Date    Date     `json:"date"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type Artist struct {
	Name string `json:"name"`
}

type Medium struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	ID string `json:"id"`
}

type AcoustMatch struct {
	AlbumID     string
	TrackID     string
	AlbumName   string
	ArtistName  string
	ReleaseDate Date
}

type AcoustID struct {
	apiKey      string
	log         *slog.Logger
	lastReqTime time.Time
}

func (a *AcoustID) GetMusicBrainzAlbumID(filename string) ([]AcoustMatch, error) {
	a.log.Info(fmt.Sprintf("analyzing file '%s'", filename))
	fpCalc, err := a.fpCalc(filename)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("client", a.apiKey)
	params.Set("meta", "releases compress tracks")
	params.Set("duration", fmt.Sprint(math.Round(float64(fpCalc.Duration))))
	params.Set("fingerprint", fpCalc.Fingerprint)

	u := fmt.Sprintf("%s?%s", apiUrl, params.Encode())

	if time.Now().Before(a.lastReqTime.Add(minReqTimeout)) {
		a.log.Debug(fmt.Sprintf("sleeping for %d ms to avoid rate limit", minReqTimeout-time.Since(a.lastReqTime)))
		time.Sleep(minReqTimeout - time.Since(a.lastReqTime))
	}

	resp, err := http.Get(u)
	a.lastReqTime = time.Now()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Results) == 0 {
		// return nil, errors.New("could not find a match")
		return nil, nil
	}

	if response.Results[0].Score < scoreLimit {
		// return nil, errors.New("could not find a good enough match")
		return nil, nil
	}

	matches := []AcoustMatch{}

	for _, r := range response.Results[0].Releases {
		if len(r.Mediums) == 0 || len(r.Mediums[0].Tracks) == 0 || len(r.Artists) == 0 {
			continue
		}
		matches = append(matches, AcoustMatch{
			AlbumID:     r.ID,
			TrackID:     r.Mediums[0].Tracks[0].ID,
			AlbumName:   r.Title,
			ArtistName:  r.Artists[0].Name,
			ReleaseDate: r.Date,
		})
	}

	if len(matches) == 0 {
		// return nil, errors.New("no releases found")
		return nil, nil
	}

	return matches, nil
}

func (a AcoustID) fpCalc(filename string) (fpCalcResp, error) {
	fpCalcRes, err := exec.Command("fpcalc", "-json", filename).Output()
	if err != nil {
		return fpCalcResp{}, err
	}

	var fpCalc fpCalcResp
	if err = json.Unmarshal(fpCalcRes, &fpCalc); err != nil {
		return fpCalcResp{}, err
	}

	return fpCalc, nil
}

func New(apiKey string, slogHandler slog.Handler) (AcoustID, error) {
	_, err := exec.LookPath("fpcalc")
	if err != nil {
		return AcoustID{}, errors.New("'fpcalc' command not found on the computer")
	}

	return AcoustID{
		apiKey: apiKey,
		log:    slog.New(slogHandler).With("service", "acousticID"),
	}, nil
}
