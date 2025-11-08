package wiki

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/tidwall/gjson"
)

const (
	wikipediaUrlTemplate = "https://%s/api/rest_v1/page/summary/%s"
)

type Wikipedia struct {
	Summary string
}

func (w Wiki) getWikipedia(ctx context.Context, wikipediaUrl string) (Wikipedia, error) {
	u, err := url.Parse(wikipediaUrl)
	if err != nil {
		return Wikipedia{}, err
	}

	wikiID := path.Base(u.Path)
	domain := u.Hostname()

	req, err := http.NewRequest("GET", fmt.Sprintf(wikipediaUrlTemplate, domain, wikiID), nil)
	if err != nil {
		return Wikipedia{}, err
	}

	resp, err := w.do(ctx, req)
	if err != nil {
		return Wikipedia{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Wikipedia{}, fmt.Errorf("returned status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Wikipedia{}, err
	}

	wp := Wikipedia{}

	res := gjson.ParseBytes(bodyBytes)

	if res.Get("type").String() != "standard" {
		return Wikipedia{}, fmt.Errorf("cannot use an article with '%s' type", res.Get("type").String())
	}

	wp.Summary = res.Get("extract").String()

	return wp, nil
}
