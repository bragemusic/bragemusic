package musicbrainz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	baseUrl       = "https://musicbrainz.org/ws/2"
	minReqTimeout = 1100 * time.Millisecond
)

var (
	preferredMediaOrder = []string{"cd", "digital"}
	excludeMediaTypes   = []string{"dvd"}
)

type MusicBrainz struct {
	lastReqTime time.Time
	log         *slog.Logger
}

func (m *MusicBrainz) GetArtist(ctx context.Context, artistID string) (ArtistResponse, error) {
	params := url.Values{}
	params.Set("inc", "aliases tags genres url-rels")
	params.Set("fmt", "json")

	u := fmt.Sprintf("%s/artist/%s?%s", baseUrl, artistID, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return ArtistResponse{}, err
	}

	resp, err := m.do(ctx, req)
	if err != nil {
		return ArtistResponse{}, err
	}
	defer resp.Body.Close()

	var artistData ArtistResponse
	if err = json.NewDecoder(resp.Body).Decode(&artistData); err != nil {
		return ArtistResponse{}, err
	}

	return artistData, nil
}

func (m *MusicBrainz) GetAlbumFromNames(ctx context.Context, artist, album string) (*Release, error) {
	params := url.Values{}
	params.Set("query", fmt.Sprintf(`release:"%s" AND artist:"%s" AND status:Official`, album, artist))
	params.Set("fmt", "json")

	u := fmt.Sprintf("%s/release/?%s", baseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releaseData ReleaseSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&releaseData); err != nil {
		return nil, err
	}

	releases := releaseData.Releases

	for _, preferredMedia := range preferredMediaOrder {
		filteredReleases := m.filterMedia(releaseData.Releases, preferredMedia)
		if len(filteredReleases) > 0 {
			releases = filteredReleases
			break
		}
	}

	score := 0
	releaseDate := time.Now()
	found := false
	var release Release
	for _, r := range releases {
		if r.Score >= score && r.Date.Before(releaseDate) && r.Date.After(time.Time{}) {
			if found && (r.Date.Day() == 1 && r.Date.Month() == 1) && (releaseDate.Day() != 1 && releaseDate.Month() != 1) {
				continue
			}
			release = r
			score = r.Score
			releaseDate = r.Date.Time
			found = true
		}
	}

	if !found {
		// return Release{}, fmt.Errorf("could not find release for '%s' by '%s'", album, artist)
		return nil, nil
	}

	media, err := m.getReleaseDetails(ctx, release)
	if err != nil {
		return nil, err
	}

	release.Media = media

	return &release, nil
}

func (m *MusicBrainz) GetAlbum(ctx context.Context, releaseID string) (Release, error) {
	params := url.Values{}
	params.Set("inc", `recordings media artist-credits url-rels`)
	params.Set("fmt", "json")

	u := fmt.Sprintf("%s/release/%s?%s", baseUrl, releaseID, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Release{}, err
	}

	resp, err := m.do(ctx, req)
	if err != nil {
		return Release{}, err
	}
	defer resp.Body.Close()

	var releaseData Release
	if err = json.NewDecoder(resp.Body).Decode(&releaseData); err != nil {
		return Release{}, err
	}

	return releaseData, nil
}

func (m *MusicBrainz) getReleaseDetails(ctx context.Context, release Release) ([]Media, error) {
	params := url.Values{}
	params.Set("inc", `recordings media`)
	params.Set("fmt", "json")

	u := fmt.Sprintf("%s/release/%s?%s", baseUrl, release.ID, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releaseData Release
	if err = json.NewDecoder(resp.Body).Decode(&releaseData); err != nil {
		return nil, err
	}

	var media []Media
	for _, m := range releaseData.Media {
		if slices.Contains(excludeMediaTypes, strings.ToLower(m.Format)) {
			continue
		}
		media = append(media, m)
	}

	return media, nil
}

func (m MusicBrainz) DownloadCoverArt(ctx context.Context, albumMbID, albumID string, outputFolder string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://coverartarchive.org/release/%s/front", albumMbID), nil)
	if err != nil {
		return err
	}

	resp, err := m.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	out, err := os.Create(filepath.Join(outputFolder, albumID+".jpg"))
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

func (m MusicBrainz) filterMedia(releases []Release, mediaType string) []Release {
	filteredReleases := []Release{}
	for _, r := range releases {
		if len(r.Media) == 0 {
			continue
		}
		if strings.Contains(strings.ToLower(r.Media[0].Format), strings.ToLower(mediaType)) {
			filteredReleases = append(filteredReleases, r)
		}
	}
	return filteredReleases
}

func (m *MusicBrainz) do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	if time.Now().Before(m.lastReqTime.Add(minReqTimeout)) {
		m.log.Debug(fmt.Sprintf("sleeping for %d ms to avoid rate limit", minReqTimeout-time.Since(m.lastReqTime)))
		time.Sleep(minReqTimeout - time.Since(m.lastReqTime))
	}

	backoff := 3 * time.Second

	for i := range 3 {
		if i > 0 {
			m.log.DebugContext(ctx, "trying again", "attempt", i+1)
		}
		client := &http.Client{}
		resp, err = client.Do(req)
		m.lastReqTime = time.Now()
		if err != nil {
			m.log.DebugContext(ctx, "backing off", "error", err.Error())
			time.Sleep(time.Duration(i+1) * backoff)
		} else {
			return resp, nil
		}
	}

	return nil, errors.New("could not request musicbrainz")
}

func New(slogHandler slog.Handler) MusicBrainz {
	return MusicBrainz{
		log: slog.New(slogHandler).With("service", "musicbrainz"),
	}
}
