package audioreader

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
)

type AudioReader interface {
	ReadMediafile(ctx context.Context, mf types.MediaFile) (types.MediaStream, error)
}

func ReadOsFile(ctx context.Context, filename string) (types.MediaStream, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return fileReader{
		File: f,
	}, nil
}

type LocalAudioReader struct {
	musicdirPath string
}

func (r LocalAudioReader) ReadMediafile(ctx context.Context, mf types.MediaFile) (types.MediaStream, error) {
	filename := filepath.Join(r.musicdirPath, mf.Filename())
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return fileReader{
		File: f,
	}, nil
}

type ServerAudioReader struct {
	sc  *serverclient.ServerClient
	log *slog.Logger
}

func (r ServerAudioReader) ReadMediafile(ctx context.Context, mf types.MediaFile) (types.MediaStream, error) {
	// FIXME: DONT HARDCODE
	chunkSize := 1 * 1024 * 1024 // 4MB
	// sr := streamingReader{mediafile: mf, sc: r.sc, chunkSize: chunkSize, log: r.log}
	sr := newStreamingReader(r.sc, mf, chunkSize, r.log)

	// err := sr.LoadChunk()
	// if err != nil {
	// 	return nil, err
	// }

	return sr, nil
}

func NewLocalReader(musicDir string) AudioReader {
	return LocalAudioReader{
		musicdirPath: musicDir,
	}
}

func NewServerReader(sc *serverclient.ServerClient, slogHandler slog.Handler) AudioReader {
	return ServerAudioReader{
		log: slog.New(slogHandler).With("service", "audioreader"),
		sc:  sc,
	}
}

// har fipplar jag
// Men det har ar fel tank med Readern. Vi ska nog ha en funktion pa den som tar mediafile som input och ger en io.Reader som output.
// 	nu ar nog tanket med local reader ratt
