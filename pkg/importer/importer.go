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
	"github.com/bragemusic/core/pkg/audioreader"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/filetx"
	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/sse"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
	"github.com/gofrs/uuid/v5"
)

type Config struct {
	ImportDirPath          string
	ManualImportDirPath    string
	FinishedImportsDirPath string
	MusicDirPath           string
	ImageDirPath           string
	DeleteImportsOnSuccess bool
}

type Importer struct {
	importDir       string
	manualImportDir string
	postImportDir   string
	musicDir        string
	imageDir        string
	deleteOnSuccess bool
	db              database.DatabaseFace
	mb              musicbrainz.MusicBrainz
	aid             acoustid.AcoustID
	im              imagemagick.ImageMagick
	ar              audioreader.AudioReader
	sseDispatch     sse.Dispatcher
	log             *slog.Logger
	berr            bragerr.BragErrFactory
}

func (i Importer) AddImportEntry(ctx context.Context, filename string, itype types.ImportType, userID uuid.UUID, musicbrainzID *string) error {
	ie := types.Import{
		MusicBrainzID: musicbrainzID,
		Owner:         userID,
		Filename:      filename,
		Type:          itype,
		State:         types.ImportStateNotStarted,
	}

	_, err := i.db.AddImport(ctx, ie)
	if err != nil {
		return i.berr.DatabaseError(err, types.EntityImport, nil)
	}

	if err = i.sseDispatch.Broadcast(types.SSEImporterItemsUpdated()); err != nil {
		i.log.WarnContext(ctx, "could not send updated event", "error", err.Error())
	}

	return nil
}

func (i *Importer) runImportCheck(ctx context.Context) error {
	if err := os.MkdirAll(i.postImportDir, 0o755); err != nil {
		return err
	}

	for {
		ie, found, err := i.db.GetUnclaimedImport(ctx)
		if err != nil {
			return i.berr.DatabaseError(err, types.EntityImport, nil)
		}

		if !found {
			i.log.DebugContext(ctx, "no items found for processing")
			return nil
		}

		path := filepath.Join(i.importDir, ie.Filename)

		if err = i.db.SetImportState(ctx, ie.ID, types.ImportStateRunning); err != nil {
			return i.berr.DatabaseError(err, types.EntityImport, &ie.ID)
		}

		if err = i.sseDispatch.Broadcast(types.SSEImporterItemsUpdated()); err != nil {
			i.log.WarnContext(ctx, "could not send updated event", "error", err.Error())
		}

		switch ie.Type {
		case types.ImportTypeAlbum:
			err = i.importAlbum(ctx, path, ie.Owner, ie.MusicBrainzID)
		case types.ImportTypeTrack:
			err = i.importTrack(ctx, path)
		}

		if err != nil {
			if dberr := i.db.SetImportState(ctx, ie.ID, types.ImportStateError); dberr != nil {
				return i.berr.DatabaseError(dberr, types.EntityImport, &ie.ID)
			}

			if err = i.sseDispatch.Broadcast(types.SSEImporterItemsUpdated()); err != nil {
				i.log.WarnContext(ctx, "could not send updated event", "error", err.Error())
			}

			i.log.ErrorContext(ctx, "import failed", "type", ie.Type, "filename", ie.Filename, "error", err.Error())
			continue
		}

		if err = i.db.SetImportState(ctx, ie.ID, types.ImportStateFinished); err != nil {
			return i.berr.DatabaseError(err, types.EntityImport, &ie.ID)
		}

		if err = i.sseDispatch.Broadcast(types.SSEImporterItemsUpdated()); err != nil {
			i.log.WarnContext(ctx, "could not send updated event", "error", err.Error())
		}

		if !i.deleteOnSuccess {
			err = i.copyFile(ctx, path, filepath.Join(i.postImportDir, ie.Filename))
			if err != nil {
				return err
			}
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}

	}
}

func (i *Importer) runManualDirImportCheck(ctx context.Context) error {
	if err := os.MkdirAll(i.importDir, 0o755); err != nil {
		return err
	}

	return filepath.Walk(i.manualImportDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				i.log.ErrorContext(ctx, "skipping file due to error", "error", err.Error())
			}

			switch strings.ToLower(filepath.Ext(path)) {
			case ".zip":
				i.log.InfoContext(ctx, "manual import file found", "filename", path)
				if err := i.AddImportEntry(ctx, filepath.Base(path), types.ImportTypeAlbum, uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")), nil); err != nil {
					return err
				}

				err = i.copyFile(ctx, path, filepath.Join(i.importDir, filepath.Base(path)))
				if err != nil {
					return err
				}

				err = os.Remove(path)
				if err != nil {
					return err
				}
			default:
				return nil
			}
			return nil
		})
}

func (i *Importer) importAlbum(ctx context.Context, filename string, userID uuid.UUID, mbID *string) error {
	i.log.InfoContext(ctx, "parsing album", "filename", filename)

	tempFolder, err := os.MkdirTemp(os.TempDir(), "brage-album")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempFolder)

	if err = i.unzipMusicFiles(ctx, filename, tempFolder); err != nil {
		return err
	}

	if err = i.importAlbumFiles(ctx, tempFolder, userID, mbID); err != nil {
		return err
	}

	return nil
}

func (i Importer) unzipMusicFiles(ctx context.Context, filename, targetDir string) error {
	i.log.InfoContext(ctx, "unzipping file", "filename", filename)

	archive, err := zip.OpenReader(filename)
	if err != nil {
		return err
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
	tempFolder, err := os.MkdirTemp(os.TempDir(), "brage-img")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempFolder)

	if err = i.downloadAlbumCoverImage(ctx, album, mdPictures, tempFolder); err != nil {
		return err
	}

	imgFilename := filepath.Join(tempFolder, album.ID.String()+".jpg")

	outputFolder := filepath.Join(i.imageDir, "albums", album.ID.String())

	if err = i.im.ResizeAll(ctx, imgFilename, outputFolder); err != nil {
		return err
	}

	return nil
}

func (i Importer) downloadAlbumCoverImage(ctx context.Context, album types.Album, mdPictures []*tag.Picture, outputFolder string) error {
	if album.MusicBrainzID != nil {
		i.log.InfoContext(ctx, "downloading album cover from MusicBrainz", "album", album.Name)

		err := i.mb.DownloadCoverArt(ctx, *album.MusicBrainzID, album.ID.String(), outputFolder)
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

		imgFilename := filepath.Join(outputFolder, fmt.Sprintf("%s.%s", album.ID, pic.Ext))
		err := utils.SaveID3Image(ctx, *pic, imgFilename)
		if err != nil {
			i.log.WarnContext(ctx, "could not get image from ID3", "error", err.Error())
			continue
		}
		return nil
	}

	return fmt.Errorf("could not get album cover for album '%s'", album.ID)
}

func (i Importer) importAlbumFiles(ctx context.Context, folder string, userID uuid.UUID, mbID *string) error {
	ftx, err := filetx.Begin(ctx)
	if err != nil {
		return err
	}
	defer ftx.Rollback()

	mediaFiles, err := i.importMediaFiles(ctx, &ftx, folder, userID)
	if err != nil {
		return err
	}

	albumAnalysis, err := i.analyzeAlbum(ctx, mediaFiles, mbID)
	if err != nil {
		if errors.Is(err, ErrAlbumMbIDNotFound) {
			i.log.WarnContext(ctx, "no album musicbrainz ID found, using ID3")
		} else {
			return err
		}
	} else {
		i.log.InfoContext(ctx, "album musicbrainz ID found", "mbID", albumAnalysis.AlbumID)
	}

	existingAlbum, err := i.getExistingAlbum(ctx, albumAnalysis)
	if err != nil {
		if !errors.Is(err, ErrAlbumNotFound) {
			return err
		}
	}

	artist, err := i.generateArtist(ctx, albumAnalysis, userID)
	if err != nil {
		return err
	}

	tx, err := i.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, mf := range mediaFiles {
		exists, err := tx.MediaFileExists(ctx, mf.ID)
		if err != nil {
			return err
		}

		if !exists {
			_, err = tx.AddMediaFile(ctx, mf, userID)
			if err != nil {
				return err
			}
		}
	}

	albumTracks, err := i.addMultipleTracks(ctx, tx, albumAnalysis, existingAlbum.ID, userID)
	if err != nil {
		return err
	}

	album, err := i.addAlbum(ctx, tx, albumAnalysis, existingAlbum, userID)
	if err != nil {
		return err
	}

	if err = i.addAlbumTracks(ctx, tx, albumTracks, album.ID, userID); err != nil {
		return err
	}

	artistID, err := i.addArtist(ctx, tx, artist, albumAnalysis, userID)
	if err != nil {
		return err
	}

	if err = i.addAlbumArtists(ctx, tx, album.ID, artistID, userID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := ftx.Commit(); err != nil {
		return err
	}

	err = i.downloadAlbumCover(ctx, album, albumAnalysis.Covers)
	if err != nil {
		i.log.ErrorContext(ctx, "could not download album cover", "id", album.ID, "error", err.Error())
	}

	return nil
}

func (i *Importer) Run(ctx context.Context) error {
	i.log.InfoContext(ctx, "starting import check")

	err := i.runManualDirImportCheck(ctx)
	if err != nil {
		return err
	}

	err = i.runImportCheck(ctx)
	if err != nil {
		return err
	}

	i.log.InfoContext(ctx, "import check done")
	return nil
}

func New(cfg Config, db database.DatabaseFace, sd sse.Dispatcher, mb musicbrainz.MusicBrainz, aid acoustid.AcoustID, im imagemagick.ImageMagick, slogHandler slog.Handler) Importer {
	return Importer{
		importDir:       cfg.ImportDirPath,
		manualImportDir: cfg.ManualImportDirPath,
		musicDir:        cfg.MusicDirPath,
		imageDir:        cfg.ImageDirPath,
		db:              db,
		mb:              mb,
		aid:             aid,
		im:              im,
		ar:              audioreader.NewLocalReader(cfg.MusicDirPath),
		sseDispatch:     sd,
		log:             slog.New(slogHandler).With("service", "importer"),
		postImportDir:   cfg.FinishedImportsDirPath,
		deleteOnSuccess: cfg.DeleteImportsOnSuccess,
		berr:            bragerr.NewFactory("importer"),
	}
}
