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
)

const playCountCountFrac = 0.75

var ErrFileNotFound = errors.New("file not found")

type Config struct {
	PlayerName   string
	MusicDirPath string
}

type CallbackType string

type AudioPlayer struct {
	ai                            audiointerface.AudioInterface
	mp                            mpris.Mpris
	playCtx                       PlayContext
	currentFile                   *os.File
	progressTicker                *time.Ticker
	currentPlayCtxChangeCallbacks []func(PlayContext)
	pausePlayCallbacks            []func(isPlaying bool)
	progressCallbacks             []func(ms int64)
	playCountCallbacks            []func(trackID string)
	errCallback                   func(context.Context, error)
	musicDirPath                  string
	playCountReported             bool
	log                           *slog.Logger
}

func (a *AudioPlayer) RegisterErrorCallback(f func(context.Context, error)) {
	a.errCallback = f
}

func (a *AudioPlayer) RegisterPlatContextChangeCallback(f func(PlayContext)) {
	a.currentPlayCtxChangeCallbacks = append(a.currentPlayCtxChangeCallbacks, f)
}

func (a *AudioPlayer) RegisterPlayPauseCallback(f func(isPlaying bool)) {
	a.pausePlayCallbacks = append(a.pausePlayCallbacks, f)
}

func (a *AudioPlayer) RegisterProgressCallback(f func(ms int64)) {
	a.progressCallbacks = append(a.progressCallbacks, f)
}

func (a *AudioPlayer) RegisterPlayCountCallback(f func(trackID string)) {
	a.playCountCallbacks = append(a.playCountCallbacks, f)
}

func (a *AudioPlayer) LoadAndStartTracks(ctx context.Context, playCtx PlayContext) (err error) {
	fmt.Println(playCtx.Shuffle, playCtx.Repeat)
	if playCtx.CurrentTrackIdx < 0 || playCtx.CurrentTrackIdx >= len(playCtx.Tracks) {
		return errors.New("startTrackIndex must be between 0 and len of tracks")
	}

	if err = a.closeCurrentFile(ctx); err != nil {
		return err
	}

	a.playCtx = playCtx
	a.playCtx.CurrentTrack = &playCtx.Tracks[playCtx.CurrentTrackIdx]

	return a.startTrack(ctx)
}

func (a *AudioPlayer) SetRepeat(ctx context.Context, r RepeatType) {
	a.playCtx.Repeat = r

	for _, f := range a.currentPlayCtxChangeCallbacks {
		f(a.playCtx)
	}
}

func (a *AudioPlayer) SetShuffle(ctx context.Context, s bool) {
	a.playCtx.Shuffle = s

	for _, f := range a.currentPlayCtxChangeCallbacks {
		f(a.playCtx)
	}
}

func (a *AudioPlayer) NextTrack(ctx context.Context) (err error) {
	a.log.DebugContext(ctx, "next track")

	var cidx int
	if a.playCtx.Repeat == RepeatOne {
		cidx = a.playCtx.CurrentTrackIdx
	} else {
		cidx = a.playCtx.CurrentTrackIdx + 1
	}

	if cidx >= len(a.playCtx.Tracks) {
		if a.playCtx.Repeat == RepeatAll {
			cidx = 0
		} else {
			return a.Stop(ctx)
		}
	}

	a.playCtx.CurrentTrackIdx = cidx
	a.playCtx.CurrentTrack = &a.playCtx.Tracks[a.playCtx.CurrentTrackIdx]

	for _, f := range a.currentPlayCtxChangeCallbacks {
		f(a.playCtx)
	}

	return a.startTrack(ctx)
}

func (a *AudioPlayer) PlayContext() PlayContext {
	return a.playCtx
}

func (a *AudioPlayer) currentTrack() *types.TrackDetailed {
	return a.playCtx.CurrentTrack
	// if len(a.playCtx.Tracks) == 0 {
	// 	return nil
	// }

	// return &a.playCtx.Tracks[a.playCtx.CurrentTrackIdx]
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

					if !a.playCountReported {
						totMs := a.currentTrack().MediaFile.DurationMs
						percPlayed := float32(ms) / float32(totMs)
						if percPlayed > playCountCountFrac {
							for _, f := range a.playCountCallbacks {
								f(a.currentTrack().ID)
							}
							a.playCountReported = true
						}
					}
				}
			}
		}
	}()
}

func (a *AudioPlayer) Stop(ctx context.Context) error {
	if err := a.stopPlayback(ctx); err != nil {
		return err
	}

	a.playCtx = PlayContext{
		Shuffle: a.playCtx.Shuffle,
		Repeat:  a.playCtx.Repeat,
	}

	for _, f := range a.currentPlayCtxChangeCallbacks {
		f(a.playCtx)
	}

	for _, f := range a.pausePlayCallbacks {
		f(a.ai.IsPlaying())
	}

	return nil
}

func (a *AudioPlayer) stopPlayback(ctx context.Context) error {
	a.ai.Stop()

	a.playCountReported = false

	if err := a.closeCurrentFile(ctx); err != nil {
		return err
	}

	return nil
}

func (a *AudioPlayer) startTrack(ctx context.Context) (err error) {
	if err = a.stopPlayback(ctx); err != nil {
		return err
	}

	if a.currentTrack().MediaFile == nil {
		return ErrFileNotFound
	}

	trackFilePath := filepath.Join(a.musicDirPath, a.currentTrack().MediaFile.Filename())

	a.currentFile, err = os.Open(trackFilePath)
	if err != nil {
		return err
	}

	af, err := files.ParseAudioFile(a.currentFile, a.currentTrack().MediaFile.Codec)
	if err != nil {
		return err
	}
	a.log.InfoContext(ctx, "start track", "title", a.currentTrack().Title, "artist", a.currentTrack().ArtistNames, "album", a.currentTrack().AlbumName)

	a.ai.StartAudioFile(ctx, af, a.NextTrack)
	a.mp.SetStatus(mpris.MprisPlaying)
	a.mp.SetTrack(a.currentTrack())

	for _, f := range a.currentPlayCtxChangeCallbacks {
		f(a.playCtx)
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
		playCtx: PlayContext{
			Shuffle: false,
			Repeat:  RepeatOff,
		},
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
