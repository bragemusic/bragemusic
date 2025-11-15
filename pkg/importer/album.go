package importer

import (
	"context"
	"os"
	"path/filepath"
)

func (i Importer) importAlbumFiles(ctx context.Context, folder string) error {
	osFiles, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	files := []string{}
	for _, f := range osFiles {
		files = append(files, filepath.Join(folder, f.Name()))
	}

	// TODO: list AcousticID for all tracks. Check against id3 if exists. check number of tracks match. Make a desicion

	return nil
}
