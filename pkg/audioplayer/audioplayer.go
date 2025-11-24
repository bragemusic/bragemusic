package audioplayer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/files"
	"github.com/bragemusic/core/pkg/mpris"
	"github.com/bragemusic/core/pkg/types"
	"github.com/dhowden/tag"
)

var ErrFileNotFound = errors.New("file not found")

type Config struct {
	PlayerName   string
	MusicDirPath string
}

type CallbackType string

type AudioPlayer struct {
	ai                          audiointerface.AudioInterface
	mp                          mpris.Mpris
	tracks                      []types.TrackEnhanced
	currentFile                 *os.File
	currentTrackIdx             int
	progressTicker              *time.Ticker
	currentTrackChangeCallbacks []func(*types.TrackEnhanced)
	pausePlayCallbacks          []func(isPlaying bool)
	progressCallbacks           []func(ms int64)
	errCallback                 func(context.Context, error)
	musicDirPath                string
	log                         *slog.Logger
}

func (a *AudioPlayer) RegisterErrorCallback(f func(context.Context, error)) {
	a.errCallback = f
}

func (a *AudioPlayer) RegisterTrackChangeCallback(f func(*types.TrackEnhanced)) {
	a.currentTrackChangeCallbacks = append(a.currentTrackChangeCallbacks, f)
}

func (a *AudioPlayer) RegisterPlayPauseCallback(f func(isPlaying bool)) {
	a.pausePlayCallbacks = append(a.pausePlayCallbacks, f)
}

func (a *AudioPlayer) RegisterProgressCallback(f func(ms int64)) {
	a.progressCallbacks = append(a.progressCallbacks, f)
}

func (a *AudioPlayer) LoadAndStartTracks(ctx context.Context, tracks []types.TrackEnhanced, startTrackIndex int) (err error) {
	if startTrackIndex < 0 || startTrackIndex >= len(tracks) {
		return errors.New("startTrackIndex must be between 0 and len of tracks")
	}

	if err = a.closeCurrentFile(ctx); err != nil {
		return err
	}

	a.tracks = tracks
	a.currentTrackIdx = startTrackIndex

	return a.startTrack(ctx)
}

func (a *AudioPlayer) NextTrack(ctx context.Context) (err error) {
	a.log.DebugContext(ctx, "next track")
	cidx := a.currentTrackIdx + 1
	if cidx >= len(a.tracks) {
		cidx = 0
	}

	a.currentTrackIdx = cidx

	for _, f := range a.currentTrackChangeCallbacks {
		f(a.CurrentTrack())
	}

	return a.startTrack(ctx)
}

func (a *AudioPlayer) CurrentTrack() *types.TrackEnhanced {
	if len(a.tracks) == 0 {
		return nil
	}

	return &a.tracks[a.currentTrackIdx]
}

func (a *AudioPlayer) Pause(ctx context.Context) {
	a.log.DebugContext(ctx, "pause")
	a.ai.Pause()
	a.mp.SetStatus(mpris.MprisPaused)
}

func (a *AudioPlayer) Play(ctx context.Context) {
	a.log.DebugContext(ctx, "play")
	a.ai.Play()
	a.mp.SetStatus(mpris.MprisPlaying)
}

func (a *AudioPlayer) PlayPause(ctx context.Context) {
	a.ai.PlayPause()
	if a.ai.IsPlaying() {
		a.mp.SetStatus(mpris.MprisPlaying)
	} else {
		a.mp.SetStatus(mpris.MprisPaused)
	}

	for _, f := range a.pausePlayCallbacks {
		f(a.ai.IsPlaying())
	}
}

func (a *AudioPlayer) Terminate(ctx context.Context) {
	a.log.DebugContext(ctx, "terminate")
	a.ai.Terminate()
}

func (a *AudioPlayer) startProgressPrinter() {
	go func() {
		for {
			select {
			case <-a.progressTicker.C:
				if a.ai.IsPlaying() {
					ms := a.ai.PlayedMS()
					for _, f := range a.progressCallbacks {
						f(ms)
					}
				}
			}
		}
	}()
}

func (a *AudioPlayer) startTrack(ctx context.Context) (err error) {
	a.ai.Stop()

	if err = a.closeCurrentFile(ctx); err != nil {
		return err
	}

	if a.CurrentTrack().FilePath == "" {
		return ErrFileNotFound
	}

	trackFilePath := filepath.Join(a.musicDirPath, a.CurrentTrack().FilePath)

	a.currentFile, err = os.Open(trackFilePath)
	if err != nil {
		return err
	}

	af, err := files.ParseAudioFile(a.currentFile, tag.FileType(*a.CurrentTrack().MimeType))
	if err != nil {
		return err
	}
	a.log.InfoContext(ctx, "start track", "title", a.CurrentTrack().Title, "artist", *a.CurrentTrack().ArtistName, "album", *a.CurrentTrack().AlbumName)

	a.ai.StartAudioFile(ctx, af, a.NextTrack)
	a.mp.SetStatus(mpris.MprisPlaying)
	a.mp.SetTrack(a.CurrentTrack())

	for _, f := range a.currentTrackChangeCallbacks {
		f(a.CurrentTrack())
	}

	for _, f := range a.pausePlayCallbacks {
		f(a.ai.IsPlaying())
	}

	return nil
}

func (a *AudioPlayer) closeCurrentFile(ctx context.Context) error {
	if a.currentFile != nil {
		a.log.DebugContext(ctx, "closing file", "file", a.currentFile.Name())
		err := a.currentFile.Close()
		if err != nil {
			fmt.Println("FIXME err", err)
			return nil
		}
	}

	return nil
}

func (a *AudioPlayer) handleError(ctx context.Context, err error) {
	a.log.ErrorContext(ctx, err.Error())
	if a.errCallback != nil {
		a.errCallback(ctx, err)
	}
}

func New(cfg Config, ai audiointerface.AudioInterface, slogHandler slog.Handler) (ap *AudioPlayer, err error) {
	ap = &AudioPlayer{
		ai:           ai,
		currentFile:  nil,
		musicDirPath: cfg.MusicDirPath,
		log:          slog.New(slogHandler).With("service", "audioplayer"),
	}

	ai.RegisterErrorCallback(ap.handleError)

	ap.mp, err = mpris.New(cfg.PlayerName, ap.Play, ap.Pause, ap.PlayPause, ap.NextTrack, ap.NextTrack)
	if err != nil {
		return nil, err
	}

	ap.progressTicker = time.NewTicker(1 * time.Second)
	ap.startProgressPrinter()

	return ap, nil
}
