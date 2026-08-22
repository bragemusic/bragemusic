package client

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bragemusic/bragemusic/internal/vars"
	"github.com/bragemusic/bragemusic/pkg/audiointerface"
	"github.com/bragemusic/bragemusic/pkg/audioplayer"
	"github.com/bragemusic/bragemusic/pkg/audioreader"
	"github.com/bragemusic/bragemusic/pkg/authclient"
	"github.com/bragemusic/bragemusic/pkg/device"
	"github.com/bragemusic/bragemusic/pkg/jobmanager"
	"github.com/bragemusic/bragemusic/pkg/playbackcontroller"
	"github.com/bragemusic/bragemusic/pkg/serverclient"
	"github.com/bragemusic/bragemusic/pkg/types"
)

type Auth struct {
	client    authclient.AuthClient
	sc        *serverclient.ServerClient
	user      *types.UserDetails
	loggedIn  bool
	tokenType *types.TokenType
	log       *slog.Logger
}

func (a *Auth) ServerClient() *serverclient.ServerClient {
	return a.sc
}

func (a *Auth) GetUser(ctx context.Context, tokenType types.TokenType) *types.UserDetails {
	if a.user != nil {
		return a.user
	}

	err := a.client.LoginCachedServerUser(ctx, tokenType)
	if err != nil {
		a.log.WarnContext(ctx, "could not log in cached user", "error", err.Error())
		return nil
	}

	user, userErr := a.sc.GetUser(ctx)
	if userErr != nil {
		a.log.WarnContext(ctx, "could not get user", "error", userErr.Error())
		if tokenType != types.TokenAPI {
			return nil
		}

		serr, ok := userErr.(serverclient.ErrStatus)
		if !ok {
			return nil
		}

		if serr.Refused && tokenType == types.TokenAPI {
			cachedUser, err := a.GetCachedUser(ctx)
			if err != nil {
				a.log.WarnContext(ctx, "could not get cached user", "error", err.Error())
				return nil
			}

			if cachedUser == nil {
				a.log.InfoContext(ctx, "no cached user found")
				return nil
			}

			user = *cachedUser

		}
	}

	a.login(ctx, user, tokenType)

	return a.user
}

func (a *Auth) GetCachedUser(ctx context.Context) (user *types.UserDetails, err error) {
	return a.client.GetCachedUser(ctx)
}

func (a *Auth) GetCachedUsers(ctx context.Context) (users []types.UserDetails, err error) {
	return a.client.GetCachedUsers(ctx)
}

func (a *Auth) login(ctx context.Context, user types.UserDetails, tokenType types.TokenType) {
	a.user = &user
	a.loggedIn = true
	a.tokenType = &tokenType

	a.log.InfoContext(ctx, "successfully logged in", "user.name", user.Username, "user.email", user.Email)
}

func (a *Auth) Logout(ctx context.Context) error {
	err := a.sc.Logout(ctx)
	if err != nil {
		return err
	}

	a.user = nil
	a.loggedIn = false

	tokenType := types.TokenFrontendShort
	if a.tokenType != nil {
		tokenType = *a.tokenType
	}

	err = a.client.LogoutServerUser(ctx, tokenType)
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) LoginServerUser(ctx context.Context, username, password string, tokenType types.TokenType) error {
	user, err := a.client.Login(ctx, username, password, tokenType)
	if err != nil {
		return err
	}

	a.login(ctx, user, tokenType)

	return nil
}

func (a *Auth) LoginToken(ctx context.Context, token string) error {
	user, err := a.client.LoginToken(ctx, token)
	if err != nil {
		return err
	}

	a.login(ctx, user, types.TokenAPI)

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

	da := device.NewAgent(slogHandler, sc, a.user.ID, types.DeviceBase{
		Name:             config.PlayerName,
		Type:             config.ClientType,
		Interface:        config.ClientInterface,
		Icon:             config.ClientIcon,
		SupportsPlayback: true,
		Platform:         vars.PLATFORM,
		Version:          vars.VERSION,
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
		a.log.ErrorContext(ctx, "could not subscribe to device events", "error", err.Error())
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
