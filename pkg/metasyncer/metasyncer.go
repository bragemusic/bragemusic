package metasyncer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/wiki"
)

type MetaSyncer struct {
	db       database.DatabaseFace
	mb       musicbrainz.MusicBrainz
	wiki     wiki.Wiki
	imageDir string
	log      *slog.Logger
}

func (m MetaSyncer) artistHasImage(artistID string) bool {
	imgFilename := filepath.Join(m.imageDir, "artists", artistID+".jpg")
	if _, err := os.Stat(imgFilename); err != nil {
		return false
	}
	return true
}

func (m MetaSyncer) getArtistMetaData(ctx context.Context, artistMbId string) (wiki.WikiData, error) {
	mbArtist, err := m.mb.GetArtist(ctx, artistMbId)
	if err != nil {
		return wiki.WikiData{}, err
	}

	wikiDataUrl := ""

	for _, rel := range mbArtist.Relations {
		if rel.Type == "wikidata" {
			wikiDataUrl = rel.URL.Resource
		}
	}

	if wikiDataUrl == "" {
		return wiki.WikiData{}, fmt.Errorf("could not get wikidata for artist MbID with '%s'", artistMbId)
	}

	wikiData, err := m.wiki.GetWikiData(ctx, wikiDataUrl)
	if err != nil {
		return wiki.WikiData{}, err
	}

	return wikiData, nil
}

func (m MetaSyncer) Sync(ctx context.Context) {
	m.log.InfoContext(ctx, "started metadata sync")
	defer func() { m.log.InfoContext(ctx, "metadata sync finsished") }()

	if err := os.MkdirAll(filepath.Join(m.imageDir, "artists"), 0o755); err != nil {
		m.log.ErrorContext(ctx, "could not create image dir", "error", err.Error())
	}

	artists, err := m.db.ListArtistsWithoutMeta(ctx)
	if err != nil {
		m.log.ErrorContext(ctx, "could not list artists", "error", err.Error())
	}
	m.log.InfoContext(ctx, fmt.Sprintf("found %d artist in need of sync", len(artists)))

	for _, a := range artists {
		wikiData, err := m.getArtistMetaData(ctx, *a.MusicBrainzID)
		if err != nil {
			m.log.ErrorContext(ctx, "could not get wiki data for artist", "error", err.Error(), "artist", a.Name)
			continue
		}

		a.Description = wikiData.Summary
		err = m.db.UpdateArtist(ctx, a)
		if err != nil {
			m.log.ErrorContext(ctx, "could not update artist", "error", err.Error(), "artist", a.Name)
			continue
		}

		m.log.InfoContext(ctx, "updated artist description", "artist", a.Name)

		if wikiData.ImageUrl != nil && !m.artistHasImage(a.ID.String()) {
			imgFilename := filepath.Join(m.imageDir, "artists", a.ID.String()+".jpg")
			if err = m.wiki.DownloadFile(ctx, *wikiData.ImageUrl, imgFilename); err != nil {
				m.log.ErrorContext(ctx, "could not download artist image", "error", err.Error(), "artist", a.Name)
				continue
			}
			m.log.InfoContext(ctx, "downloaded artist art", "artist", a.Name)
		}
	}
}

func New(imageDir string, db database.DatabaseFace, mb musicbrainz.MusicBrainz, wiki wiki.Wiki, slogHandler slog.Handler) MetaSyncer {
	return MetaSyncer{
		db:       db,
		log:      slog.New(slogHandler).With("service", "meta-syncer"),
		mb:       mb,
		wiki:     wiki,
		imageDir: imageDir,
	}
}
