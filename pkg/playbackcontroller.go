package playbackcontroller

import (
	"context"
	"log/slog"

	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/device"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type PlaybackControllerFace interface {
	audioplayer.AudioPlayerFace
	ConnectDevice(ctx context.Context, id uuid.UUID) error
	DisconnectDevice(ctx context.Context) error
}

type PlaybackController struct {
	remoteDeviceID *uuid.UUID
	localPlayer    audioplayer.AudioPlayerFace
	deviceAgent    *device.DeviceAgent
	// db             database.DatabaseFace
	sc *serverclient.ServerClient

	contextCallbacks  []func(context.Context, types.PlayContext)
	playbackCallbacks []func(context.Context, types.PlaybackState)

	log *slog.Logger
}

func (p *PlaybackController) RegisterErrorCallback(f func(context.Context, error)) {
	p.localPlayer.RegisterErrorCallback(f)
}

func (p *PlaybackController) RegisterPlayContextCallback(f func(context.Context, types.PlayContext)) {
	p.localPlayer.RegisterPlayContextCallback(f)
	p.contextCallbacks = append(p.contextCallbacks, f)
}

func (p *PlaybackController) RegisterPlaybackStateCallback(f func(context.Context, types.PlaybackState)) {
	p.localPlayer.RegisterPlaybackStateCallback(f)
	p.playbackCallbacks = append(p.playbackCallbacks, f)
}

func (p *PlaybackController) RegisterPlayCountCallback(f func(trackID uuid.UUID)) {
	p.localPlayer.RegisterPlayCountCallback(f)
}

func (p *PlaybackController) ConnectDevice(ctx context.Context, id uuid.UUID) error {
	device, err := p.deviceAgent.GetDevice(ctx, id)
	if err != nil {
		return err
	}

	if err := p.localPlayer.Stop(ctx); err != nil {
		return err
	}

	p.remoteDeviceID = &id

	ps := device.PlayerState

	if ps == nil {
		ps = &types.PlayerStateDTO{
			Playback: types.PlaybackStateDTO{
				DeviceID: id,
				PlaybackState: types.PlaybackState{
					Shuffle:     false,
					Repeat:      types.RepeatOff,
					Playing:     false,
					ProgressMS:  0,
					TrackSource: types.TrackSourceContext,
					TrackIndex:  0,
				},
			},
		}
	}

	for _, f := range p.playbackCallbacks {
		f(ctx, ps.Playback.PlaybackState)
	}

	for _, f := range p.contextCallbacks {
		f(ctx, ps.Context.PlayContext)
	}

	return nil
}

func (p *PlaybackController) DisconnectDevice(ctx context.Context) error {
	p.remoteDeviceID = nil

	ps := p.localPlayer.PlayerState()

	for _, f := range p.playbackCallbacks {
		f(ctx, ps.Playback)
	}

	for _, f := range p.contextCallbacks {
		f(ctx, ps.Context)
	}

	return nil
}

func (p *PlaybackController) PlayerState() types.PlayerState {
	if p.remoteDeviceID == nil {
		return p.localPlayer.PlayerState()
	}

	device, err := p.deviceAgent.GetDevice(context.Background(), *p.remoteDeviceID)
	if err != nil {
		p.log.Error("could not get remote device", "error", err.Error())
		return types.PlayerState{}
	}

	psdto := device.PlayerState

	if psdto == nil {
		psdto = &types.PlayerStateDTO{
			Playback: types.PlaybackStateDTO{
				DeviceID: *p.remoteDeviceID,
				PlaybackState: types.PlaybackState{
					Shuffle:     false,
					Repeat:      types.RepeatOff,
					Playing:     false,
					ProgressMS:  0,
					TrackSource: types.TrackSourceContext,
					TrackIndex:  0,
				},
			},
		}
	}

	ps := types.PlayerState{
		Playback: psdto.Playback.PlaybackState,
		Context:  psdto.Context.PlayContext,
	}

	return ps
}

func (p *PlaybackController) LoadAndStartTracks(ctx context.Context, state types.PlayerState) (err error) {
	if p.remoteDeviceID == nil {
		return p.localPlayer.LoadAndStartTracks(ctx, state)
	}

	return p.sc.DeviceSetPlayerState(ctx, *p.remoteDeviceID, state)
}

func (p *PlaybackController) NextTrack(ctx context.Context) (err error) {
	if p.remoteDeviceID == nil {
		p.localPlayer.NextTrack(ctx)
		return
	}

	return p.sc.DeviceNextTrack(ctx, *p.remoteDeviceID)
}

func (p *PlaybackController) PreviousTrack(ctx context.Context) (err error) {
	if p.remoteDeviceID == nil {
		p.localPlayer.PreviousTrack(ctx)
		return
	}

	return p.sc.DevicePreviousTrack(ctx, *p.remoteDeviceID)
}

// func (p *PlaybackController) Pause(ctx context.Context) {
// 	p.localPlayer.Pause(ctx)
// }

// func (p *PlaybackController) Play(ctx context.Context) {
// 	p.localPlayer.Play(ctx)
// }

func (p *PlaybackController) PlayPause(ctx context.Context) {
	if p.remoteDeviceID == nil {
		p.localPlayer.PlayPause(ctx)
		return
	}

	if err := p.sc.DevicePlayPause(ctx, *p.remoteDeviceID); err != nil {
		p.log.ErrorContext(ctx, "could not run command on remote device", "cmd", "play-pause")
		return
	}
}

func (p *PlaybackController) SetRepeat(ctx context.Context, r types.RepeatType) {
	if p.remoteDeviceID == nil {
		p.localPlayer.SetRepeat(ctx, r)
		return
	}

	if err := p.sc.DevicePlayerSetRepeat(ctx, *p.remoteDeviceID, r); err != nil {
		p.log.ErrorContext(ctx, "could not run command on remote device", "cmd", "set-repeat")
		return
	}
}

func (p *PlaybackController) SetShuffle(ctx context.Context, s bool) {
	if p.remoteDeviceID == nil {
		p.localPlayer.SetShuffle(ctx, s)
		return
	}

	if err := p.sc.DevicePlayerSetShuffle(ctx, *p.remoteDeviceID, s); err != nil {
		p.log.ErrorContext(ctx, "could not run command on remote device", "cmd", "set-shuffle")
		return
	}
}

func (p *PlaybackController) Stop(ctx context.Context) error {
	if p.remoteDeviceID == nil {
		return p.localPlayer.Stop(ctx)
	}

	return p.sc.DevicePlayerStop(ctx, *p.remoteDeviceID)
}

func (p *PlaybackController) AddTrackToQueue(ctx context.Context, track types.TrackDetailed) {
	p.localPlayer.AddTrackToQueue(ctx, track)
}

func (p *PlaybackController) handleRemotePlayContexts(ctx context.Context, e types.SSEvent) {
	if p.remoteDeviceID == nil {
		return
	}

	pc, err := types.DecodeEventData[types.PlayContextDTO](e)
	if err != nil {
		p.log.ErrorContext(ctx, "could not decode playcontext data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
		return
	}

	if pc.DeviceID != *p.remoteDeviceID {
		return
	}

	for _, f := range p.contextCallbacks {
		f(ctx, pc.PlayContext)
	}
}

func (p *PlaybackController) handleRemotePlaybackStates(ctx context.Context, e types.SSEvent) {
	if p.remoteDeviceID == nil {
		return
	}

	ps, err := types.DecodeEventData[types.PlaybackStateDTO](e)
	if err != nil {
		p.log.ErrorContext(ctx, "could not decode playbackstate data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
		return
	}

	if ps.DeviceID != *p.remoteDeviceID {
		return
	}

	for _, f := range p.playbackCallbacks {
		f(ctx, ps.PlaybackState)
	}
}

func (p *PlaybackController) handleDeviceDisconnection(ctx context.Context, e types.SSEvent) {
	if p.remoteDeviceID == nil {
		return
	}

	d, err := types.DecodeEventData[types.Device](e)
	if err != nil {
		p.log.ErrorContext(ctx, "could not decode device data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
		return
	}

	if d.ID == *p.remoteDeviceID {
		if err := p.DisconnectDevice(ctx); err != nil {
			p.log.ErrorContext(ctx, "could not disconnect device", "error", err.Error())
			return
		}
	}
}

func New(ap audioplayer.AudioPlayerFace, da *device.DeviceAgent, sc *serverclient.ServerClient, slogHandler slog.Handler) (PlaybackControllerFace, error) {
	// var err error

	pc := &PlaybackController{
		localPlayer: ap,
		deviceAgent: da,
		sc:          sc,
		// db:                db,
		contextCallbacks:  []func(context.Context, types.PlayContext){},
		playbackCallbacks: []func(context.Context, types.PlaybackState){},
		log:               slog.New(slogHandler).With("service", "playbackcontroller"),
	}

	da.SubscribeToEventTypes(pc.handleRemotePlayContexts, types.SSEventTypePlayerPlayContext)
	da.SubscribeToEventTypes(pc.handleRemotePlaybackStates, types.SSEventTypePlayerPlaybackState)
	da.SubscribeToEventTypes(pc.handleDeviceDisconnection, types.SSEventTypeDeviceDisconnected)
	// ap := &AudioPlayer{
	// 	ai:           ai,
	// 	ar:           ar,
	// 	currentFile:  nil,
	// 	musicDirPath: cfg.MusicDirPath,
	// 	log:          slog.New(slogHandler).With("service", "audioplayer"),
	// 	state: types.PlayerState{
	// 		Playback: types.PlaybackState{
	// 			Shuffle:     false,
	// 			Repeat:      types.RepeatOff,
	// 			TrackSource: types.TrackSourceContext,
	// 		},
	// 	},
	// }

	// ai.RegisterErrorCallback(ap.handleError)

	// playerName := "BrageMusic-" + strings.ReplaceAll(cfg.PlayerName, " ", "")
	// ap.mp, err = mpris.New(playerName, ap.Play, ap.Pause, ap.PlayPause, ap.PreviousTrack, ap.NextTrack)
	// if err != nil {
	// 	return nil, err
	// }

	// ap.progressTicker = time.NewTicker(1 * time.Second)
	// ap.startProgressPrinter()

	return pc, nil
}
