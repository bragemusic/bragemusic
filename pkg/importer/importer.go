package importer

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bragemusic/core/pkg/acoustid"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/wiki"
)

type Importer struct {
	importDir string
	musicDir  string
	db        database.DatabaseFace
	mb        musicbrainz.MusicBrainz
	aid       acoustid.AcoustID
	wiki      wiki.Wiki
	log       *slog.Logger
}

func (i *Importer) runImportCheck(ctx context.Context) error {
	return filepath.Walk(i.importDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			var f func(context.Context, string) error

			switch strings.ToLower(filepath.Ext(path)) {
			case ".flac":
				f = i.importTrack
			case ".zip":
				f = i.importAlbum
			default:
				return nil
			}

			i.log.InfoContext(ctx, "file found", "filename", path)
			err = f(ctx, path)
			if err != nil {
				i.log.ErrorContext(ctx, "could not import track", "error", err.Error())
				return err
			}

			return nil
		})
}

func (i *Importer) importAlbum(ctx context.Context, filename string) error {
	i.log.InfoContext(ctx, "parsing album", "filename", filename)

	tempFolder, err := os.MkdirTemp(os.TempDir(), "brage-album")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempFolder)

	if err = i.unzipMusicFiles(ctx, filename, tempFolder); err != nil {
		return err
	}

	if err = i.importAlbumFiles(ctx, tempFolder); err != nil {
		return err
	}

	return nil
}

func (i Importer) unzipMusicFiles(ctx context.Context, filename, targetDir string) error {
	i.log.InfoContext(ctx, "unzipping file", "filename", filename)

	archive, err := zip.OpenReader(filename)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}
		filePath := filepath.Join(targetDir, filepath.Base(f.Name))

		if !strings.HasPrefix(filePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path '%s'", filePath)
		}

		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func (i *Importer) importTrack(ctx context.Context, filename string) error {
	return errors.New("track import not implemented")
	// f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	// if err := i.trackManager.AddTrack(ctx, f); err != nil {
	// 	i.log.ErrorContext(ctx, "could not import track", "error", err.Error())
	// 	// return err
	// }
	// return nil
}

func (i *Importer) Run(ctx context.Context) {
	i.log.InfoContext(ctx, "starting import check")

	err := i.runImportCheck(ctx)
	if err != nil {
		i.log.ErrorContext(ctx, "import check finished with errors", "error", err.Error())
	}

	i.log.InfoContext(ctx, "import check done")
}

func New(importDir, musicDir string, db database.DatabaseFace, mb musicbrainz.MusicBrainz, aid acoustid.AcoustID, wiki wiki.Wiki, slogHandler slog.Handler) Importer {
	return Importer{
		importDir: importDir,
		musicDir:  musicDir,
		db:        db,
		mb:        mb,
		aid:       aid,
		wiki:      wiki,
		log:       slog.New(slogHandler).With("service", "importer"),
	}
}
