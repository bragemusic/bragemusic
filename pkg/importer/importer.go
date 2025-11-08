package importer

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bragemusic/core/pkg/trackmgr"
)

type Importer struct {
	importDir    string
	log          *slog.Logger
	trackManager *trackmgr.TrackManager
}

func (i *Importer) runImportCheck(ctx context.Context) error {
	err := filepath.Walk(i.importDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.ToLower(filepath.Ext(path)) == ".flac" {
				i.log.Info("file found", "filename", path)
				err = i.importTrack(ctx, path)
				if err != nil {
					i.log.Error("could not import track", "error", err.Error())
					return err
				}
				return nil
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (i *Importer) importTrack(ctx context.Context, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := i.trackManager.AddTrack(ctx, f); err != nil {
		i.log.ErrorContext(ctx, "could not import track", "error", err.Error())
		// return err
	}
	return nil
}

func (i *Importer) Run(ctx context.Context) {
	i.runImportCheck(ctx)

	artists, err := i.trackManager.ListArtists(ctx)
	if err != nil {
		panic(err)
	}

	// for _, a := range artists {
	// 	if a.MusicBrainzID != nil {
	// 		err = i.trackManager.GetArtistMetaData(ctx, *a.MusicBrainzID)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }

	// return

	for _, a := range artists {
		albums, err := i.trackManager.GetAlbumsByArtist(ctx, a.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println(a.Name)
		for _, al := range albums {
			fmt.Println(" ", al.Name)
			tracks, err := i.trackManager.GetTracksByAlbum(ctx, al.ID)
			if err != nil {
				panic(err)
			}
			for _, t := range tracks {
				fmt.Println("    ", *t.DiscNumber, "|", *t.TrackNumber, "-", t.Title)
			}
		}
	}
}

func New(importDir string, trackManager *trackmgr.TrackManager, slogHandler slog.Handler) Importer {
	return Importer{
		importDir:    importDir,
		log:          slog.New(slogHandler),
		trackManager: trackManager,
	}
}
