package importer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bragemusic/bragemusic/pkg/audioreader"
	"github.com/bragemusic/bragemusic/pkg/files"
	"github.com/bragemusic/bragemusic/pkg/filetx"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/bragemusic/bragemusic/pkg/utils"
	"github.com/gofrs/uuid/v5"
)

func (i Importer) importMediaFiles(ctx context.Context, ftx *filetx.FileTx, folder string, userID uuid.UUID) (mediaFiles []types.MediaFile, err error) {
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

		mf, err := i.db.GetMediaFileFromChecksum(ctx, checksum)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			i.log.InfoContext(ctx, "creating new media file")

			// f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
			// if err != nil {
			// 	return nil, err
			// }

			f, err := audioreader.ReadOsFile(ctx, filename)
			if err != nil {
				return nil, err
			}

			// FIXME: Do not hardcode Flac
			af, err := files.ParseAudioFile(f, types.CodecFlac)
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

			uid, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			mf.ID = uid

			mffp := filepath.Join(i.musicDir, fmt.Sprintf("%s.%s", uid.String(), mf.Codec))

			if err = ftx.CopyFile(filename, mffp); err != nil {
				return nil, err
			}

		} else {
			i.log.InfoContext(ctx, "media file already imported")
		}

		mf.OrgFilename = f.Name()
		mediaFiles = append(mediaFiles, mf)
	}

	return mediaFiles, nil
}
