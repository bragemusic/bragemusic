package importer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/files"
	"github.com/bragemusic/core/pkg/filetx"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
)

func (i Importer) importMediaFiles(ctx context.Context, tx database.DatabaseFace, ftx *filetx.FileTx, folder string) (mediaFiles []types.MediaFile, err error) {
	tempFiles, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, f := range tempFiles {
		filename := filepath.Join(folder, f.Name())
		checksum, err := utils.FileSHA256(filename)
		if err != nil {
			return nil, err
		}

		mf, err := tx.GetMediaFileFromChecksum(ctx, checksum)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			i.log.InfoContext(ctx, "creating new media file")

			f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
			if err != nil {
				return nil, err
			}

			// FIXME: Do not hardcode Flac
			af, err := files.ParseAudioFile(f, tag.FLAC)
			if err != nil {
				f.Close()
				return nil, err
			}
			f.Close()

			mf = types.MediaFile{
				DurationMs: af.DurationMS(),
				Bitrate:    af.Bitrate(),
				SampleRate: af.SampleRate(),
				FileSize:   af.FileSize(),
				// FIXME: Do not hardcode Flac
				Codec:    types.CodecFlac,
				Checksum: checksum,
			}

			mfId, err := tx.AddMediaFile(ctx, mf)
			if err != nil {
				return nil, err
			}

			mf.ID = mfId

			mffp := filepath.Join(i.musicDir, fmt.Sprintf("%s.%s", mfId.String(), mf.Codec))

			if err = ftx.CopyFile(filename, mffp); err != nil {
				return nil, err
			}

		} else {
			i.log.InfoContext(ctx, "media file already imported")
		}

		mediaFiles = append(mediaFiles, mf)
	}

	return mediaFiles, nil
}
