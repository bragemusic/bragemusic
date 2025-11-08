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
	wikiDataUrlFormat = "https://www.wikidata.org/wiki/Special:EntityData/%s.json"
	imageProptertyTag = "P18"
)

type WikiData struct {
	ImageUrl *string
	Summary  *string
}

func (w Wiki) GetWikiData(ctx context.Context, wikiDataUrl string) (WikiData, error) {
	u, err := url.Parse(wikiDataUrl)
	if err != nil {
		return WikiData{}, err
	}

	wikiID := path.Base(u.Path)

	req, err := http.NewRequest("GET", fmt.Sprintf(wikiDataUrlFormat, wikiID), nil)
	if err != nil {
		return WikiData{}, err
	}

	resp, err := w.do(ctx, req)
	if err != nil {
		return WikiData{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return WikiData{}, err
	}

	wd := WikiData{}

	res := gjson.ParseBytes(bodyBytes)

	imageFilename := res.Get(fmt.Sprintf("entities.%s.claims.%s.0.mainsnak.datavalue.value", wikiID, imageProptertyTag)).String()
	if imageFilename != "" {
		imageUrl := fmt.Sprintf(imageUrlFormat, imageFilename)
		wd.ImageUrl = &imageUrl
	}

	for _, lang := range w.preferredLangs {
		u := res.Get(fmt.Sprintf("entities.%s.sitelinks.%swiki.url", wikiID, lang)).String()
		if u != "" {
			wpData, err := w.getWikipedia(ctx, u)
			if err == nil {
				wd.Summary = &wpData.Summary
				break
			}
		}
	}

	return wd, nil
}
