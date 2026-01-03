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
	"github.com/bragemusic/core/pkg/filetx"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
)

type Config struct {
	ImportDirPath          string
	FinishedImportsDirPath string
	MusicDirPath           string
	ImageDirPath           string
	DeleteImportsOnSuccess bool
}

type Importer struct {
	importDir       string
	postImportDir   string
	musicDir        string
	imageDir        string
	deleteOnSuccess bool
	db              database.DatabaseFace
	mb              musicbrainz.MusicBrainz
	aid             acoustid.AcoustID
	log             *slog.Logger
}

func (i *Importer) runImportCheck(ctx context.Context) error {
	if err := os.MkdirAll(i.postImportDir, 0o755); err != nil {
		return err
	}

	return filepath.Walk(i.importDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				i.log.ErrorContext(ctx, "skipping file due to error", "error", err.Error())
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

			if !i.deleteOnSuccess {
				err = i.copyFile(ctx, path, filepath.Join(i.postImportDir, filepath.Base(path)))
				if err != nil {
					return err
				}
			}

			err = os.Remove(path)
			if err != nil {
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

		err := i.mb.DownloadCoverArt(ctx, *album.MusicBrainzID, album.ID.String(), dir)
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

func (i Importer) importAlbumFiles(ctx context.Context, folder string) error {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ftx, err := filetx.Begin(ctx)
	if err != nil {
		return err
	}
	defer ftx.Rollback()

	mediaFiles, err := i.importMediaFiles(ctx, tx, &ftx, folder)
	if err != nil {
		return err
	}

	albumAnalysis, err := i.analyzeAlbum(ctx, mediaFiles)
	if err != nil {
		if errors.Is(err, ErrAlbumMbIDNotFound) {
			i.log.WarnContext(ctx, "no album musicbrainz ID found, using ID3")
		} else {
			return err
		}
	} else {
		i.log.InfoContext(ctx, "album musicbrainz ID found", "mbID", albumAnalysis.AlbumID)
	}

	existingAlbum, err := i.getExistingAlbum(ctx, tx, albumAnalysis)
	if err != nil {
		if !errors.Is(err, ErrAlbumNotFound) {
			return err
		}
	}

	albumTracks, err := i.addMultipleTracks(ctx, tx, albumAnalysis, existingAlbum.ID)
	if err != nil {
		return err
	}

	album, err := i.addAlbum(ctx, tx, albumAnalysis, existingAlbum)
	if err != nil {
		return err
	}

	if err = i.addAlbumTracks(ctx, tx, albumTracks, album.ID); err != nil {
		return err
	}

	artistID, err := i.addArtist(ctx, tx, albumAnalysis)
	if err != nil {
		return err
	}

	if err = i.addAlbumArtists(ctx, tx, album.ID, artistID); err != nil {
		return err
	}

	err = i.downloadAlbumCover(ctx, album, albumAnalysis.Covers)
	if err != nil {
		i.log.ErrorContext(ctx, "could not download album cover", "id", album.ID, "error", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := ftx.Commit(); err != nil {
		return err
	}

	return nil
}

func (i *Importer) Run(ctx context.Context) {
	i.log.InfoContext(ctx, "starting import check")

	err := i.runImportCheck(ctx)
	if err != nil {
		i.log.ErrorContext(ctx, "import check finished with errors", "error", err.Error())
	}

	i.log.InfoContext(ctx, "import check done")
}

func New(cfg Config, db database.DatabaseFace, mb musicbrainz.MusicBrainz, aid acoustid.AcoustID, slogHandler slog.Handler) Importer {
	return Importer{
		importDir:       cfg.ImportDirPath,
		musicDir:        cfg.MusicDirPath,
		imageDir:        cfg.ImageDirPath,
		db:              db,
		mb:              mb,
		aid:             aid,
		log:             slog.New(slogHandler).With("service", "importer"),
		postImportDir:   cfg.FinishedImportsDirPath,
		deleteOnSuccess: cfg.DeleteImportsOnSuccess,
	}
}
