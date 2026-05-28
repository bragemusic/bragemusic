package client

import (
	"context"
	"errors"
	"log/slog"

	playbackcontroller "github.com/bragemusic/core/pkg"
	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/audioreader"
	"github.com/bragemusic/core/pkg/authclient"
	"github.com/bragemusic/core/pkg/device"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type Auth struct {
	client   authclient.AuthClient
	sc       *serverclient.ServerClient
	user     *types.UserDetails
	loggedIn bool
	log      *slog.Logger
}

func (a *Auth) ServerClient() *serverclient.ServerClient {
	return a.sc
}

func (a *Auth) GetUser() *types.UserDetails {
	return a.user
}

func (a *Auth) GetCachedUsers(ctx context.Context) (users []types.UserDetails, err error) {
	return a.client.GetCachedUsers(ctx)
}

func (a *Auth) login(ctx context.Context, user types.UserDetails) {
	a.user = &user
	a.loggedIn = true

	a.log.InfoContext(ctx, "successfully logged in", "user.name", user.Username, "user.email", user.Email)
}

func (a *Auth) LoginServerUser(ctx context.Context, username, password string, longLivedToken bool) error {
	user, err := a.client.Login(ctx, username, password, longLivedToken)
	if err != nil {
		return err
	}

	a.login(ctx, user)

	return nil
}

func (a *Auth) LoginLocalUser(ctx context.Context, userID uuid.UUID) error {
	user, err := a.client.LoginLocalUser(ctx, userID, false)
	if err != nil {
		return err
	}

	a.login(ctx, user)

	return nil
}

func (a *Auth) LoginToken(ctx context.Context, token string) error {
	user, err := a.client.LoginToken(ctx, token)
	if err != nil {
		return err
	}

	a.login(ctx, user)

	return nil
}

func (a *Auth) ServerStatus(ctx context.Context) (types.ServerApiInfo, error) {
	sai, err := a.sc.CheckStatus(ctx)
	if err != nil {
		sai = types.ServerApiInfo{
			ServerInfo: types.ServerInfo{
				Status: types.HealthzUnavailable,
			},
		}
	}
	return sai, nil
}

func (a *Auth) NewClient(ctx context.Context, config Config, slogHandler slog.Handler) (ClientFace, error) {
	if a.user == nil || !a.loggedIn {
		return nil, errors.New("not logged in")
	}

	ctx, cancelFunc := context.WithCancel(ctx)

	var cf clientFace
	var ar audioreader.AudioReader
	var err error

	sc := a.client.ServerClient()

	jm := jobmanager.New(slogHandler)

	switch config.ClientType {
	case types.DeviceTypeStreaming:
		cf, err = newStreamingClient(ctx, config, &jm, sc, *a.user, slogHandler)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		ar = audioreader.NewServerReader(sc, slogHandler)
	case types.DeviceTypeSync:
		cf, err = newSyncClient(ctx, config, &jm, sc, *a.user, slogHandler)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		ar = audioreader.NewLocalReader(config.MusicDirPath)
	default:
		cancelFunc()
		return nil, errors.New("type not implemented")
	}

	id, err := NewIdentity("lucas test")
	if err != nil {
		cancelFunc()
		return nil, err
	}

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	apCfg := audioplayer.Config{
		PlayerName:   config.PlayerName,
		MusicDirPath: config.MusicDirPath,
	}

	ap, err := audioplayer.New(ctx, apCfg, pa, ar, slogHandler)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	// FIXME: need to fix platform and version
	da := device.NewAgent(slogHandler, sc, a.user.ID, types.DeviceBase{
		Name:             config.PlayerName,
		Type:             config.ClientType,
		Interface:        config.ClientInterface,
		Icon:             config.ClientIcon,
		SupportsPlayback: true,
		Platform:         "linux",
		Version:          "1.2.0",
	},
		config.StateFilePath,
	)

	// FIXME: Only for visual clients
	pc := playbackcontroller.New(ap, da, sc, slogHandler)

	c := &Client{
		clientFace:             cf,
		Identity:               &id,
		PlaybackControllerFace: pc,
		sc:                     sc,
		log:                    slog.New(slogHandler).With("service", "client"),
		DeviceAgent:            da,
		activeRemoteDevice:     nil,
		contextCallbacks:       []func(context.Context, types.PlayContext){},
		playbackCallbacks:      []func(context.Context, types.PlaybackState){},
		closeFunc:              cancelFunc,
	}

	err = c.SubscribeDeviceEvents(ctx)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	ap.RegisterPlayCountCallback(c.updatePlayCount)
	ap.RegisterPlayContextCallback(c.handlePlayContextCallbacks)
	ap.RegisterPlaybackStateCallback(c.handlePlaybackStateCallbacks)

	da.SubscribeToEventCategory(c.handlePlayerEvents, "player")
	da.SubscribeToEventCategory(c.handleServerEvent, "importer")
	da.SubscribeToEventCategory(c.handleServerEvent, "server")

	jm.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobAuthClientServerStatus,
		CronExpr: "*/10 * * * * *",
		Run:      c.UpdateServerStatus,
	})

	return c, nil
}

func NewFromToken(ctx context.Context, token string, config Config, slogHandler slog.Handler) (ClientFace, error) {
	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	a := Auth{
		client: authclient.New(&sc, slogHandler),
		log:    slog.New(slogHandler).With("service", "auth-client"),
	}

	err := a.LoginToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return a.NewClient(ctx, config, slogHandler)
}

func NewAuthClient(ctx context.Context, config Config, slogHandler slog.Handler) Auth {
	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	return Auth{
		client: authclient.New(&sc, slogHandler),
		sc:     &sc,
		log:    slog.New(slogHandler).With("service", "auth-client"),
	}
}
