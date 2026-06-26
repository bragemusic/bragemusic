// Package client provides the main client-side API for interacting with
// the system.
//
// It exposes a unified interface that combines authentication, synchronization,
// audio playback control, metadata access, and background job management.
// The client package acts as the primary entry point for applications that
// need to communicate with and operate against the server and local storage.
package client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bragemusic/core/internal/vars"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/device"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/playbackcontroller"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/sse"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type Config struct {
	ConfigPath      string
	ImagePath       string
	MusicDirPath    string
	PlayerName      string
	ServerBaseURL   string
	ClientType      types.DeviceType
	ClientInterface types.DeviceInterface
	ClientIcon      types.DeviceIcon
	StateFilePath   *string
}

type clientSync struct {
	*syncer.Syncer
	*mediamanager.MediaManager
	*jobmanager.JobManager

	sc *serverclient.ServerClient

	config  Config
	log     *slog.Logger
	berr    bragerr.BragErrFactory
	dbClose func() error

	eventCallbacks []func(types.ClientEvent, any)
	user           types.UserDetails
}

type clientStreaming struct {
	*jobmanager.JobManager
	*syncer.NoSync

	*serverclient.ServerClient

	config Config
	log    *slog.Logger
	berr   bragerr.BragErrFactory

	eventCallbacks []func(types.ClientEvent, any)
	user           types.UserDetails
}

type Client struct {
	clientFace
	*Identity
	playbackcontroller.PlaybackControllerFace

	*device.DeviceAgent

	activeRemoteDevice *uuid.UUID

	contextCallbacks  []func(context.Context, types.PlayContext)
	playbackCallbacks []func(context.Context, types.PlaybackState)

	user types.UserDetails

	serverStatus             types.ServerApiInfo
	serverAvailableCallbacks []func(types.ServerApiInfo)

	updateServerStatusCallbacks []func(context.Context)

	serverEventCallbacks []sse.EventHandler

	closeFunc context.CancelFunc

	sc  *serverclient.ServerClient
	log *slog.Logger
}

func (c *Client) RegisterUpdateServerStatusCallback(f func(context.Context)) {
	c.updateServerStatusCallbacks = append(c.updateServerStatusCallbacks, f)
}

func (c *Client) RegisterServerAvailabilityCallback(f func(types.ServerApiInfo)) {
	c.serverAvailableCallbacks = append(c.serverAvailableCallbacks, f)
}

func (c *Client) UpdateServerStatusCallbacks(ctx context.Context) {
	for _, f := range c.updateServerStatusCallbacks {
		f(ctx)
	}
}

func (c *Client) handleServerEvent(ctx context.Context, event types.SSEvent) {
	for _, f := range c.serverEventCallbacks {
		f(ctx, event)
	}
}

func (c *Client) Close(ctx context.Context) error {
	err := c.PlaybackControllerFace.Close(ctx)
	if err != nil {
		return err
	}

	c.closeFunc()

	return nil
}

func (c *Client) SubscribeToClientEvents(handler sse.EventHandler) {
	c.serverEventCallbacks = append(c.serverEventCallbacks, handler)

	c.DeviceAgent.SubscribeToClientEvents(handler)
}

func (c *Client) StartPlayerWithAlbum(ctx context.Context, albumID uuid.UUID, trackNumber int) error {
	tracks, err := c.ListTracksDetailedByAlbum(ctx, albumID)
	if err != nil {
		return err
	}

	pState := types.PlayerState{
		Playback: types.PlaybackState{
			TrackIndex: trackNumber,
		},
		Context: types.PlayContext{
			Type:   types.PlayContextAlbum,
			RefID:  albumID,
			Tracks: tracks,
		},
	}

	err = c.PlaybackControllerFace.LoadAndStartTracks(ctx, pState)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "albumID", albumID.String(), "trackNumber", trackNumber)

	return nil
}

func (c *Client) StartPlayerWithLikedTracks(ctx context.Context, trackNumber int) error {
	tracks, err := c.ListLikedTracks(ctx)
	if err != nil {
		return err
	}

	pState := types.PlayerState{
		Playback: types.PlaybackState{
			TrackIndex: trackNumber,
		},
		Context: types.PlayContext{
			Type:   types.PlayContextLikedTracks,
			RefID:  uuid.Nil,
			Tracks: tracks,
		},
	}

	err = c.PlaybackControllerFace.LoadAndStartTracks(ctx, pState)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "type", "liked tracks", "trackNumber", trackNumber)

	return nil
}

func (c *Client) StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, sortBy database.SortBy, sortOrder database.SortOrder) error {
	tracks, err := c.ListPlaylistTracks(ctx, playlistID, sortBy, sortOrder)
	if err != nil {
		return err
	}

	pState := types.PlayerState{
		Playback: types.PlaybackState{
			TrackIndex: trackNumber,
		},
		Context: types.PlayContext{
			Type:   types.PlayContextPlaylist,
			RefID:  playlistID,
			Tracks: tracks,
		},
	}

	err = c.PlaybackControllerFace.LoadAndStartTracks(ctx, pState)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "playlistID", playlistID.String(), "trackNumber", trackNumber)

	return nil
}

func (c *Client) updatePlayCount(trackID uuid.UUID) {
	// FIXME: Do we need context here?
	ctx := context.TODO()

	err := c.AddPlayCount(ctx, trackID)
	if err != nil {
		c.log.ErrorContext(ctx, "could not add play count", "error", err.Error())
		return
	}
	c.log.DebugContext(ctx, "added play count", "track_id", trackID)

	// c.emitEvent(types.ClientEventEntitiesUpdated, nil)
}

func (c *Client) AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error {
	track, err := c.GetTrackDetailed(ctx, trackID, albumID)
	if err != nil {
		return err
	}

	c.PlaybackControllerFace.AddTrackToQueue(ctx, track)
	return nil
}

func (c *Client) handlePlayerEvents(ctx context.Context, e types.SSEvent) {
	switch e.Type {
	case types.SSEventTypePlayerAddToQueue:
		track, err := types.DecodeEventData[types.TrackDetailed](e)
		if err != nil {
			c.log.ErrorContext(ctx, "could not decode playerstate data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		c.PlaybackControllerFace.AddTrackToQueue(ctx, track)
	case types.SSEventTypePlayerNextTrack:
		if err := c.NextTrack(ctx); err != nil {
			c.log.ErrorContext(ctx, "could not execute remote command", "command", e.Type, "error", err.Error())
		}
	case types.SSEventTypePlayerPlayPause:
		c.PlayPause(ctx)
	case types.SSEventTypePlayerPreviousTrack:
		if err := c.PreviousTrack(ctx); err != nil {
			c.log.ErrorContext(ctx, "could not execute remote command", "command", e.Type, "error", err.Error())
		}
	case types.SSEventTypePlayerSetRepeat:
		rt, err := types.DecodeEventData[types.RepeatType](e)
		if err != nil {
			c.log.ErrorContext(ctx, "could not decode playerstate data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}
		c.SetRepeat(ctx, rt)
	case types.SSEventTypePlayerSetShuffle:
		s, err := types.DecodeEventData[bool](e)
		if err != nil {
			c.log.ErrorContext(ctx, "could not decode playerstate data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}
		c.SetShuffle(ctx, s)
	case types.SSEventTypePlayerSetState:
		ps, err := types.DecodeEventData[types.PlayerState](e)
		if err != nil {
			c.log.ErrorContext(ctx, "could not decode playerstate data in event", "event.type", e.Type, "event.id", e.ID.String(), "event.data", e.Data)
			return
		}

		if err = c.LoadAndStartTracks(ctx, ps); err != nil {
			c.log.ErrorContext(ctx, "could not load state", "error", err.Error())
			return
		}
	case types.SSEventTypePlayerStop:
		if err := c.Stop(ctx); err != nil {
			c.log.ErrorContext(ctx, "could not execute remote command", "command", e.Type, "error", err.Error())
		}
	}
}

func (c *Client) handlePlayContextCallbacks(ctx context.Context, pc types.PlayContext) {
	err := c.sc.UpdatePlayContext(ctx, pc)
	if err != nil {
		c.log.ErrorContext(ctx, "could not send playcontext to server", "error", err.Error())
	}
}

func (c *Client) handlePlaybackStateCallbacks(ctx context.Context, ps types.PlaybackState) {
	err := c.sc.UpdatePlaybackState(ctx, ps)
	if err != nil {
		c.log.ErrorContext(ctx, "could not send playback state to server", "error", err.Error())
	}
}

func (c *Client) UpdateServerStatus(ctx context.Context) error {
	h, err := c.sc.CheckStatus(ctx)
	if err != nil {
		h = types.ServerApiInfo{
			ServerInfo: types.ServerInfo{
				Status: types.HealthzUnavailable,
			},
		}
	}

	c.serverStatus = h

	for _, f := range c.serverAvailableCallbacks {
		f(h)
	}
	if err != nil {
		return err
	}

	if h.Application != vars.SERVER_APP_NAME {
		return fmt.Errorf("expected server application name differs. '%s' != expected '%s'", h.Application, vars.SERVER_APP_NAME)
	}

	if h.Status == types.HealthzUnavailable {
		return fmt.Errorf("server status is not running. '%s'", h.Status)
	}
	return nil
}

func (c *Client) ServerStatus(ctx context.Context) (types.ServerApiInfo, error) {
	if err := c.UpdateServerStatus(ctx); err != nil {
		return types.ServerApiInfo{}, err
	}
	return c.serverStatus, nil
}

// func New(ctx context.Context, config Config, slogHandler slog.Handler) (ClientFace, error) {
// 	var cf clientFace
// 	var ar audioreader.AudioReader
// 	var err error

// 	sc := serverclient.New(config.ServerBaseURL, slogHandler)

// 	switch config.ClientType {
// 	case types.DeviceTypeStreaming:
// 		cf, err = newStreamingClient(ctx, config, &sc, slogHandler)
// 		if err != nil {
// 			return nil, err
// 		}
// 		ar = audioreader.NewServerReader(&sc, slogHandler)
// 	case types.DeviceTypeSync:
// 		cf, err = newSyncClient(ctx, config, &sc, slogHandler)
// 		if err != nil {
// 			return nil, err
// 		}
// 		ar = audioreader.NewLocalReader(config.MusicDirPath)
// 	default:
// 		return nil, errors.New("type not implemented")
// 	}

// 	id, err := NewIdentity("lucas test")
// 	if err != nil {
// 		return nil, err
// 	}

// 	pa, err := audiointerface.NewPortAudio(slogHandler)
// 	if err != nil {
// 		return nil, err
// 	}

// 	apCfg := audioplayer.Config{
// 		PlayerName:   config.PlayerName,
// 		MusicDirPath: config.MusicDirPath,
// 	}

// 	ap, err := audioplayer.New(apCfg, pa, ar, slogHandler)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// FIXME: Needs to make another client structure. With authed starting the real clients
// 	da := device.NewAgent(slogHandler, &sc, uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")), types.DeviceBase{
// 		Name:             config.PlayerName,
// 		Type:             config.ClientType,
// 		Interface:        config.ClientInterface,
// 		Icon:             config.ClientIcon,
// 		SupportsPlayback: true,
// 		Platform:         "linux",
// 		Version:          "1.2.0",
// 	},
// 		config.StateFilePath,
// 	)

// 	// FIXME: Only for visual clients
// 	pc := playbackcontroller.New(ap, da, &sc, slogHandler)

// 	c := &Client{
// 		clientFace:             cf,
// 		Identity:               &id,
// 		PlaybackControllerFace: pc,
// 		sc:                     &sc,
// 		log:                    slog.New(slogHandler).With("service", "client"),
// 		DeviceAgent:            da,
// 		activeRemoteDevice:     nil,
// 		contextCallbacks:       []func(context.Context, types.PlayContext){},
// 		playbackCallbacks:      []func(context.Context, types.PlaybackState){},
// 	}

// 	// if err := da.SubscribeDeviceEvents(ctx); err != nil {
// 	// 	return nil, err
// 	// }
// 	ap.RegisterPlayCountCallback(c.updatePlayCount)
// 	ap.RegisterPlayContextCallback(c.handlePlayContextCallbacks)
// 	ap.RegisterPlaybackStateCallback(c.handlePlaybackStateCallbacks)

// 	da.SubscribeToEventCategory(c.handlePlayerEvents, "player")

// 	return c, nil
// }
