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
	"github.com/bragemusic/core/pkg/metasyncer"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/bragemusic/core/pkg/wiki"
	"github.com/dhowden/tag"
)

type Importer struct {
	importDir string
	musicDir  string
	imageDir  string
	db        database.DatabaseFace
	ms        *metasyncer.MetaSyncer
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

		if strings.ToLower(filepath.Ext(f.Name)) != ".flac" {
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

func (i Importer) copyFile(ctx context.Context, from, to string) error {
	i.log.DebugContext(ctx, "copying file", "src", from, "dst", to)

	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		return err
	}

	src, err := os.OpenFile(from, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func (i Importer) downloadAlbumCover(ctx context.Context, album types.Album, mdPictures []*tag.Picture) error {
	dir := filepath.Join(i.imageDir, "albums")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if album.MusicBrainzID != nil {
		i.log.InfoContext(ctx, "downloading album cover from MusicBrainz", "album", album.Name)

		err := i.mb.DownloadCoverArt(ctx, *album.MusicBrainzID, album.ID, dir)
		if err == nil {
			return nil
		}

		i.log.WarnContext(ctx, "could not get album cover from MusicBrainz")
	}

	i.log.InfoContext(ctx, "grabbing from ID3 instead")
	for _, pic := range mdPictures {
		if pic == nil {
			continue
		}

		imgFilename := filepath.Join(dir, fmt.Sprintf("%s.%s", album.ID, pic.Ext))
		err := utils.SaveID3Image(ctx, *pic, imgFilename)
		if err != nil {
			i.log.WarnContext(ctx, "could not get image from ID3", "error", err.Error())
			continue
		}
		return nil
	}

	return fmt.Errorf("could not get album cover for album '%s'", album.ID)
}

func (i *Importer) Run(ctx context.Context) {
	i.log.InfoContext(ctx, "starting import check")

	err := i.runImportCheck(ctx)
	if err != nil {
		i.log.ErrorContext(ctx, "import check finished with errors", "error", err.Error())
	}

	i.ms.Sync(ctx)

	i.log.InfoContext(ctx, "import check done")
}

func New(importDir, musicDir, imageDir string, ms *metasyncer.MetaSyncer, db database.DatabaseFace, mb musicbrainz.MusicBrainz, aid acoustid.AcoustID, wiki wiki.Wiki, slogHandler slog.Handler) Importer {
	return Importer{
		importDir: importDir,
		musicDir:  musicDir,
		imageDir:  imageDir,
		db:        db,
		mb:        mb,
		aid:       aid,
		wiki:      wiki,
		log:       slog.New(slogHandler).With("service", "importer"),
		ms:        ms,
	}
}
