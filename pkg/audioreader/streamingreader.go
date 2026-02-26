package audioreader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"syscall"
	"time"

	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
)

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}

type streamingReader struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	mediafile types.MediaFile
	sc        *serverclient.ServerClient

	chunkSize int

	chunks chan []byte
	errCh  chan error

	current []byte
	closed  bool

	log *slog.Logger
}

func (s *streamingReader) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true

	s.log.Debug("closing stream")
	s.ctxCancel()

	s.chunks = nil
	s.current = nil

	return nil
}

func (s *streamingReader) Read(p []byte) (int, error) {
	if len(s.current) == 0 {
		select {
		case chunk, ok := <-s.chunks:
			if !ok {
				select {
				case err := <-s.errCh:
					return 0, err
				default:
					return 0, io.EOF
				}
			}
			s.current = chunk

		case err := <-s.errCh:
			return 0, err

		}
	}

	n := copy(p, s.current)
	s.current = s.current[n:]

	return n, nil
}

func (s *streamingReader) Seek(offset int64, whence int) (int64, error) {
	// FIXME: Implement
	return 0, nil
}

func (s streamingReader) Size() (int64, error) {
	return s.mediafile.FileSize, nil
}

func (s *streamingReader) downloadLoop() {
	defer close(s.chunks)

	start := 0
	refusedErrCount := 0
	refusedErrLimit := 10
	backoffSec := 3

	for {
		end := start + s.chunkSize - 1

		buf := make([]byte, 0, s.chunkSize)

		s.log.Debug("loading next chunk", "start", start, "end", end)
		final, err := s.sc.DownloadMediaFilePart(
			s.ctx,
			s.mediafile.ID,
			start,
			end,
			writerFunc(func(p []byte) (int, error) {
				buf = append(buf, p...)
				return len(p), nil
			}),
		)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				refusedErrCount += 1
				if refusedErrCount >= refusedErrLimit {
					s.errCh <- fmt.Errorf("server has been unreachable for %d seconds, killing stream", backoffSec*refusedErrCount)
					s.Close()
					return
				}
				time.Sleep(time.Duration(backoffSec) * time.Second)
				continue
			} else {
				s.errCh <- err
				return
			}
		}

		select {
		case s.chunks <- buf:
		case <-s.ctx.Done():
			s.log.Debug("stopping downloading of chunks")
			return
		}

		if final {
			s.log.Debug("all chunks loaded")
			return
		}

		start = end + 1
	}
}

func newStreamingReader(sc *serverclient.ServerClient, mf types.MediaFile, chunkSize int, log *slog.Logger) *streamingReader {
	ctx, cancel := context.WithCancel(context.Background())

	s := &streamingReader{
		mediafile: mf,
		sc:        sc,
		chunkSize: chunkSize,
		chunks:    make(chan []byte, 4), // buffer 4 chunks
		errCh:     make(chan error, 1),
		ctx:       ctx,
		ctxCancel: cancel,
		log:       log,
	}

	go s.downloadLoop()

	return s
}
