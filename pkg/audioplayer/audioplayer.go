package audioplayer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioreader"
	"github.com/bragemusic/core/pkg/files"
	"github.com/bragemusic/core/pkg/mpris"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

const playCountCountFrac = 0.75

var ErrFileNotFound = errors.New("file not found")

type Config struct {
	PlayerName   string
	MusicDirPath string
}

type CallbackType string

type AudioPlayer struct {
	ai                 audiointerface.AudioInterface
	ar                 audioreader.AudioReader
	mp                 mpris.Mpris
	state              types.PlayerState
	currentFile        types.MediaStream
	progressTicker     *time.Ticker
	playCountCallbacks []func(trackID uuid.UUID)
	contextCallbacks   []func(context.Context, types.PlayContext)
	playbackCallbacks  []func(context.Context, types.PlaybackState)
	errCallback        func(context.Context, error)
	musicDirPath       string
	playCountReported  bool
	log                *slog.Logger
}

func (a *AudioPlayer) RegisterErrorCallback(f func(context.Context, error)) {
	a.errCallback = f
}

func (a *AudioPlayer) RegisterPlayContextCallback(f func(context.Context, types.PlayContext)) {
	a.contextCallbacks = append(a.contextCallbacks, f)
}

func (a *AudioPlayer) RegisterPlaybackStateCallback(f func(context.Context, types.PlaybackState)) {
	a.playbackCallbacks = append(a.playbackCallbacks, f)
}

func (a *AudioPlayer) RegisterPlayCountCallback(f func(trackID uuid.UUID)) {
	a.playCountCallbacks = append(a.playCountCallbacks, f)
}

func (a *AudioPlayer) sendStateCallback(ctx context.Context) {
	a.sendContextCallback(ctx)
	a.sendPlaybackCallback(ctx)
}

func (a *AudioPlayer) sendContextCallback(ctx context.Context) {
	for _, f := range a.contextCallbacks {
		f(ctx, a.state.Context)
	}
}

func (a *AudioPlayer) sendPlaybackCallback(ctx context.Context) {
	a.state.Playback.ProgressMS = a.ai.PlayedMS()
	a.state.Playback.UpdatedAt = time.Now()

	for _, f := range a.playbackCallbacks {
		f(ctx, a.state.Playback)
	}
}

func (a *AudioPlayer) PlayerState() types.PlayerState {
	return a.state
}

func (a *AudioPlayer) LoadAndStartTracks(ctx context.Context, state types.PlayerState) (err error) {
	if state.Playback.TrackIndex < 0 || state.Playback.TrackIndex >= len(state.Context.Tracks) {
		return errors.New("startTrackIndex must be between 0 and len of tracks")
	}

	state.Playback.Repeat = a.state.Playback.Repeat
	state.Playback.Shuffle = a.state.Playback.Shuffle
	state.Context.Queue = a.state.Context.Queue
	state.Playback.TrackSource = types.TrackSourceContext

	state.RebuildTrackOrder()

	// if err = a.closeCurrentFile(ctx); err != nil {
	// 	return err
	// }

	a.state = state

	if err := a.startTrack(ctx); err != nil {
		return err
	}

	a.sendStateCallback(ctx)

	return nil
}

func (a *AudioPlayer) SetRepeat(ctx context.Context, r types.RepeatType) {
	a.state.Playback.Repeat = r
	a.sendPlaybackCallback(ctx)
}

func (a *AudioPlayer) SetShuffle(ctx context.Context, s bool) {
	a.state.Playback.Shuffle = s
	a.state.RebuildTrackOrder()

	a.sendStateCallback(ctx)
}

func (a *AudioPlayer) NextTrack(ctx context.Context) (err error) {
	a.log.DebugContext(ctx, "next track")

	contextUpdated, stop := a.state.NextTrack()

	if contextUpdated {
		a.sendContextCallback(ctx)
	}

	if stop {
		return a.Stop(ctx)
	}

	return a.startTrack(ctx)
}

func (a *AudioPlayer) PreviousTrack(ctx context.Context) (err error) {
	a.log.DebugContext(ctx, "previous track")

	stop := a.state.PreviousTrack()
	if stop {
		return a.Stop(ctx)
	}

	return a.startTrack(ctx)
}

func (a *AudioPlayer) Pause(ctx context.Context) {
	a.log.DebugContext(ctx, "pause")

	a.state.Playback.Playing = false
	a.sendPlaybackCallback(ctx)

	a.ai.Pause()
	a.mp.SetStatus(mpris.MprisPaused)
}

func (a *AudioPlayer) Play(ctx context.Context) {
	a.log.DebugContext(ctx, "play")

	a.state.Playback.Playing = true
	a.sendPlaybackCallback(ctx)

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

	a.state.Playback.Playing = a.ai.IsPlaying()
	a.sendPlaybackCallback(ctx)
}

func (a *AudioPlayer) startProgressPrinter() {
	go func() {
		for {
			select {
			case <-a.progressTicker.C:
				if a.ai.IsPlaying() {
					ms := a.ai.PlayedMS()

					if !a.playCountReported {
						ct, err := a.state.CurrentTrack()
						if err != nil {
							a.log.Warn("could not check for play count", "error", err.Error())
							continue
						}
						totMs := ct.MediaFile.DurationMs
						percPlayed := float32(ms) / float32(totMs)
						if percPlayed > playCountCountFrac {
							for _, f := range a.playCountCallbacks {
								f(ct.ID)
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

	a.state = types.PlayerState{
		Playback: types.PlaybackState{
			Shuffle:     a.state.Playback.Shuffle,
			Repeat:      a.state.Playback.Repeat,
			Playing:     false,
			ProgressMS:  0,
			TrackSource: types.TrackSourceContext,
			TrackIndex:  0,
		},
	}

	a.sendStateCallback(ctx)

	return nil
}

func (a *AudioPlayer) AddTrackToQueue(ctx context.Context, track types.TrackDetailed) {
	a.state.Context.Queue = append(a.state.Context.Queue, track)

	a.sendStateCallback(ctx)

	a.log.InfoContext(ctx, "added track to queue", "name", track.Title, "album", track.AlbumName, "artist", track.ArtistNames)
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

	cT, err := a.state.CurrentTrack()
	if err != nil {
		return err
	}

	if cT.MediaFile == nil {
		return ErrFileNotFound
	}

	a.currentFile, err = a.ar.ReadMediafile(ctx, *cT.MediaFile)
	if err != nil {
		return err
	}

	af, err := files.ParseAudioFile(a.currentFile, cT.MediaFile.Codec)
	if err != nil {
		return err
	}
	a.log.InfoContext(ctx, "start track", "title", cT.Title, "artist", cT.ArtistNames, "album", cT.AlbumName)

	a.ai.StartAudioFile(ctx, af, a.NextTrack)

	a.state.Playback.Playing = true

	// FIXME: MP should consume state aswell
	a.mp.SetStatus(mpris.MprisPlaying)
	a.mp.SetTrack(&cT)

	a.sendPlaybackCallback(ctx)

	return nil
}

func (a *AudioPlayer) closeCurrentFile(ctx context.Context) error {
	if a.currentFile != nil {
		a.log.DebugContext(ctx, "closing file")
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

func New(cfg Config, ai audiointerface.AudioInterface, ar audioreader.AudioReader, slogHandler slog.Handler) (AudioPlayerFace, error) {
	var err error

	ap := &AudioPlayer{
		ai:           ai,
		ar:           ar,
		currentFile:  nil,
		musicDirPath: cfg.MusicDirPath,
		log:          slog.New(slogHandler).With("service", "audioplayer"),
		state: types.PlayerState{
			Playback: types.PlaybackState{
				Shuffle:     false,
				Repeat:      types.RepeatOff,
				TrackSource: types.TrackSourceContext,
			},
		},
	}

	ai.RegisterErrorCallback(ap.handleError)

	playerName := "BrageMusic-" + strings.ReplaceAll(cfg.PlayerName, " ", "")
	ap.mp, err = mpris.New(playerName, ap.Play, ap.Pause, ap.PlayPause, ap.PreviousTrack, ap.NextTrack)
	if err != nil {
		return nil, err
	}

	ap.progressTicker = time.NewTicker(1 * time.Second)
	ap.startProgressPrinter()

	return ap, nil
}
