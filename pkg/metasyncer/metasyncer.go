package metasyncer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/imagemagick"
	"github.com/bragemusic/bragemusic/pkg/musicbrainz"
	"github.com/bragemusic/bragemusic/pkg/wiki"
	"github.com/gofrs/uuid/v5"
)

type MetaSyncer struct {
	db             database.DatabaseFace
	mb             musicbrainz.MusicBrainz
	wiki           wiki.Wiki
	im             imagemagick.ImageMagick
	imageDir       string
	log            *slog.Logger
	internalUserID uuid.UUID
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

func (m MetaSyncer) Sync(ctx context.Context) error {
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
		err = m.db.UpdateArtist(ctx, a, m.internalUserID)
		if err != nil {
			m.log.ErrorContext(ctx, "could not update artist", "error", err.Error(), "artist", a.Name)
			continue
		}

		m.log.InfoContext(ctx, "updated artist description", "artist", a.Name)

		if wikiData.ImageUrl != nil && !m.artistHasImage(a.ID.String()) {
			if err := m.downloadArtistImage(ctx, *wikiData.ImageUrl, a.ID); err != nil {
				m.log.ErrorContext(ctx, "could not download artist image", "error", err.Error(), "artist", a.Name)
				continue
			}
			m.log.InfoContext(ctx, "downloaded artist art", "artist", a.Name)
		}
	}
	return nil
}

func (m MetaSyncer) downloadArtistImage(ctx context.Context, imageUrl string, artistID uuid.UUID) error {
	tempFolder, err := os.MkdirTemp(os.TempDir(), "brage-img")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempFolder)

	imgFilename := filepath.Join(tempFolder, artistID.String()+".jpg")
	if err = m.wiki.DownloadFile(ctx, imageUrl, imgFilename); err != nil {
		return err
	}

	outputFolder := filepath.Join(m.imageDir, "artists", artistID.String())

	if err = m.im.ResizeAll(ctx, imgFilename, outputFolder); err != nil {
		return err
	}

	return nil
}

func New(imageDir string, db database.DatabaseFace, mb musicbrainz.MusicBrainz, wiki wiki.Wiki, im imagemagick.ImageMagick, internalUserID uuid.UUID, slogHandler slog.Handler) MetaSyncer {
	return MetaSyncer{
		db:             db,
		log:            slog.New(slogHandler).With("service", "meta-syncer"),
		mb:             mb,
		wiki:           wiki,
		im:             im,
		imageDir:       imageDir,
		internalUserID: internalUserID,
	}
}
