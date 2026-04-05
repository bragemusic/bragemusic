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
	"errors"
	"log/slog"

	playbackcontroller "github.com/bragemusic/core/pkg"
	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/audioreader"
	"github.com/bragemusic/core/pkg/authclient"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/device"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
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

type IdentityFace interface {
	ClientID() uuid.UUID
	ClientName() string
	ClientType() types.DeviceType
}

// SyncFace defines functionality for synchronizing server data
// with the local client storage.
type SyncFace interface {
	// RegisterSyncInProgressCallback registers a callback that is invoked
	// whenever the sync state changes. The callback receives true when a
	// sync starts and false when it completes.
	RegisterSyncInProgressCallback(f func(bool))

	// SupportsSync reports whether the current client configuration
	// supports synchronization with a remote server.
	SupportsSync() bool

	// Sync synchronizes server data to the local client.
	// It blocks until the operation completes or the context is canceled.
	Sync(ctx context.Context) error
}

// AuthFace defines authentication and user session related functionality,
// including local login, server login, and server availability tracking.
type AuthFace interface {
	// RegisterUserCallback registers a callback that is invoked whenever
	// the active user changes. The callback receives nil if no user is logged in.
	RegisterUserCallback(f func(*types.UserDetails))

	// LoginLocalUser logs in a locally cached user by ID without contacting the server.
	LoginLocalUser(ctx context.Context, userID uuid.UUID) error

	// LogoutLocalUser logs out the currently active local user.
	LogoutLocalUser(ctx context.Context)

	// LoginCachedServerUser logs in a previously authenticated server user
	// using a cached identity and password. If longLivedToken is true,
	// a persistent authentication token is requested.
	LoginCachedServerUser(ctx context.Context, password string, longLivedToken bool) error

	// GetUser returns the currently authenticated user, or nil if no user is logged in.
	GetUser() *types.UserDetails

	// RegisterServerAvailabilityCallback registers a callback that is invoked
	// whenever the server availability or API information changes.
	RegisterServerAvailabilityCallback(f func(types.ServerApiInfo))

	// GetCachedUsers returns users previously cached on this client.
	GetCachedUsers(ctx context.Context) (users []types.UserDetails, err error)

	// Login authenticates a user against the server using username and password.
	// If longLivedToken is true, a persistent authentication token is requested.
	Login(ctx context.Context, username, password string, longLivedToken bool) (types.UserDetails, error)

	LoginToken(ctx context.Context, token string) (types.UserDetails, error)

	// LogoutServerUser logs out the currently authenticated server user
	// and invalidates any associated server session.
	LogoutServerUser(ctx context.Context) error

	// ServerStatus retrieves the current server API information and availability state.
	ServerStatus(ctx context.Context) (types.ServerApiInfo, error)

	// RemoveToken deletes the token specified by the ID from the server
	RemoveToken(ctx context.Context, id uuid.UUID) error
}

// AudioPlayerFace defines audio playback control and playback state
// observation functionality.
type AudioPlayerFace interface {
	// StartPlayerWithAlbum starts playback of an album beginning at the given track number.
	StartPlayerWithAlbum(ctx context.Context, albumID uuid.UUID, trackNumber int) error

	// StartPlayerWithPlaylist starts playback of a playlist beginning at the given track number.
	// Tracks are ordered according to sortBy and sortOrder.
	StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, sortBy database.SortBy, sortOrder database.SortOrder) error

	// StartPlayerWithLikedTracks starts playback of an liked tracks beginning at the given track number.
	StartPlayerWithLikedTracks(ctx context.Context, trackNumber int) error

	// AddTrackToQueue adds a track to the current playback queue.
	AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error

	// RegisterPlayContextCallback registers a callback that is invoked
	// whenever the play context (album, playlist, queue, etc.) changes.
	RegisterPlayContextCallback(f func(context.Context, types.PlayContext))

	// RegisterPlaybackStateCallback registers a callback that is invoked
	// whenever the playback state (playing, shuffle, repeat, etc.) changes.
	RegisterPlaybackStateCallback(f func(context.Context, types.PlaybackState))

	// NextTrack skips to the next track in the current play context.
	NextTrack(ctx context.Context) (err error)

	// PlayPause toggles playback between playing and paused states.
	PlayPause(ctx context.Context)

	// PreviousTrack skips to the previous track in the current play context.
	PreviousTrack(ctx context.Context) (err error)

	// SetRepeat sets the repeat mode for playback.
	SetRepeat(ctx context.Context, r types.RepeatType)

	// SetShuffle enables or disables shuffle mode for playback.
	SetShuffle(ctx context.Context, s bool)

	// PlayerState returns the current player state.
	PlayerState() types.PlayerState

	ConnectDevice(ctx context.Context, id uuid.UUID) error
	DisconnectDevice(ctx context.Context) error
}

// MetadataFace defines access to and modification of music library metadata,
// including artists, albums, tracks, playlists, users, and search functionality.
type MetadataFace interface {
	// CountArtists returns the total number of artists.
	CountArtists(ctx context.Context) (int, error)

	// GetArtist returns metadata for a specific artist.
	GetArtist(ctx context.Context, artistID uuid.UUID) (types.Artist, error)

	// GetArtistTopTracks returns the top tracks for a specific artist.
	GetArtistTopTracks(ctx context.Context, artistID uuid.UUID) ([]types.TrackDetailed, error)

	// ListArtists returns artists ordered by the provided sorting options.
	ListArtists(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.ArtistDetailed, error)

	// UpdateArtist updates metadata for an artist.
	UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error

	// UploadArtistImage uploads or replaces the image associated with an artist.
	UploadArtistImage(ctx context.Context, artistID uuid.UUID, img serverclient.FileUpload) error

	// CountAlbums returns the total number of albums.
	CountAlbums(ctx context.Context) (int, error)

	// GetAlbumDetailed returns detailed metadata for a specific album.
	GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (types.AlbumDetailed, error)

	// ListAlbumsByArtist returns albums for a given artist ordered by the provided sorting options.
	ListAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error)

	// ListAlbums returns albums ordered by the provided sorting options.
	ListAlbums(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error)

	// UpdateAlbum updates metadata for an album.
	UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error

	// UploadAlbumImage uploads or replaces the image associated with an album.
	UploadAlbumImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error

	// CountTracks returns the total number of tracks.
	CountTracks(ctx context.Context) (int, error)

	// GetTrackDetailed returns the detailed version of the wanted track
	GetTrackDetailed(ctx context.Context, trackID, albumID uuid.UUID) (track types.TrackDetailed, err error)

	// ListTracksDetailedByAlbum returns detailed track metadata for a given album.
	ListTracksDetailedByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.TrackDetailed, error)

	// RateTrack sets the rating value for a specific track.
	RateTrack(ctx context.Context, trackID uuid.UUID, value int) error

	// UpdateTrack updates metadata for a track.
	UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error

	// AddPlaylist creates a new playlist.
	AddPlaylist(ctx context.Context, playlist types.Playlist) error

	// AddPlaylistTrack adds a track to a playlist.
	AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error

	// CountPlaylists returns the total number of playlists.
	CountPlaylists(ctx context.Context) (int, error)

	// CountPlaylistTracks returns the number of tracks in a playlist.
	CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error)

	// DeletePlaylist removes a playlist.
	DeletePlaylist(ctx context.Context, id uuid.UUID) error

	// DeletePlaylistTrack removes a track from a playlist.
	DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error

	// GetPlaylist returns metadata for a specific playlist.
	GetPlaylist(ctx context.Context, id uuid.UUID) (types.Playlist, error)

	// ListPlaylists returns playlists, optionally including public playlists,
	// ordered by the provided sorting options.
	ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error)

	// ListPlaylistTracks returns tracks in a playlist ordered by the provided sorting options.
	ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error)

	// UpdatePlaylist updates metadata for a playlist.
	UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error

	// UploadPlaylistImage uploads or replaces the image associated with a playlist.
	UploadPlaylistImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error

	// ListEntityEvents returns metadata-related entity events.
	ListEntityEvents(ctx context.Context) ([]types.EntityEvent, error)

	// SearchFull performs a full-text search across supported entities
	// and returns matching search items.
	SearchFull(ctx context.Context, searchTerm string) (si []types.SearchItem, err error)

	// ImportAlbum imports an album from a file, optionally using a MusicBrainz ID
	// to enrich metadata.
	ImportAlbum(ctx context.Context, filename string, musicbrainzID *string) error

	// AddPlayCount increments the play count for a track for a specific user.
	AddPlayCount(ctx context.Context, trackID uuid.UUID) error

	// LikeTrack adds a like to a track for a specific user. If the track already has a like an error is returned.
	LikeTrack(ctx context.Context, trackID uuid.UUID) error

	// UnlikeTrack removes a like to a track for a specific user. If the track does not have a like an error is returned.
	UnlikeTrack(ctx context.Context, trackID uuid.UUID) error

	// ListLikedTracks returns the tracks liked by the authenticated user.
	ListLikedTracks(ctx context.Context) ([]types.TrackDetailed, error)

	// CountLikedTracks returns the total number of tracks liked by the authenticated user.
	CountLikedTracks(ctx context.Context) (cnt int, err error)
}

type AccessFace interface {
	// CreateUser creates a new user on the server.
	CreateUser(ctx context.Context, email, username, password string, roles []types.UserRole) error

	// DeleteUser removes the selected user from the server. This action cannot be reversed.
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// ListUsers returns all users known to the system.
	ListUsers(ctx context.Context, includeMachineUsers bool) ([]types.UserDetails, error)

	// ListUserRoles returns a list of available user roles
	ListUserRoles(ctx context.Context) ([]types.UserRole, error)

	// UpdateUser updates the selected user's information. If password is not nil, it will be changed.
	UpdateUser(ctx context.Context, userID uuid.UUID, email, username string, password *string, roles []types.UserRole) error
}

// JobManagerFace defines background job execution and scheduling functionality.
type JobManagerFace interface {
	// StartScheduler starts the background job scheduler.
	// It typically runs until the context is canceled.
	StartScheduler(ctx context.Context)

	// RunJob executes a specific job type immediately.
	RunJob(ctx context.Context, jobType types.JobType) error
}

type DeviceManagerFace interface {
	ListDevices(ctx context.Context) (devices []types.DeviceDetailed, err error)
	SubscribeToEventTypes(handler sse.EventHandler, eventType ...types.SSEventType)
	DeleteDevice(ctx context.Context, deviceID uuid.UUID) error
	DeleteDeviceToken(ctx context.Context, deviceID uuid.UUID) error
	DeleteDeviceAndToken(ctx context.Context, deviceID uuid.UUID) error
}

// ClientFace defines the high-level client interface that aggregates
// synchronization, authentication, playback, metadata, and job management
// functionality.
//
// It represents the main entry point for interacting with the system from
// a client perspective.
type clientFace interface {
	// RegisterEventCallback registers a callback that is invoked whenever
	// a client-level evant is emitted. Messages may represent events,
	// warnings, or informational notifications originating from the client.
	RegisterEventCallback(f func(types.ClientEvent, any))

	SyncFace
	AuthFace
	MetadataFace
	JobManagerFace
	AccessFace
}

type ClientFace interface {
	SubscribeToClientEvents(handler sse.EventHandler)
	clientFace
	IdentityFace
	AudioPlayerFace
	DeviceManagerFace

	// hit ska åtminstone auth och audioplayer flyttas. Det är delat mellan sync och streaming. Typ iaf. Får lösa så att auth är det.
	// 	Sen ska clientface generarea en clientSync/Stream när det loggas in en användare, så slipper vi pekar på user o en massa skit
}

type clientSync struct {
	*syncer.Syncer
	authclient.AuthClient
	// *audioplayer.AudioPlayer
	*mediamanager.MediaManager
	*jobmanager.JobManager

	sc *serverclient.ServerClient

	config  Config
	log     *slog.Logger
	berr    bragerr.BragErrFactory
	dbClose func() error

	eventCallbacks []func(types.ClientEvent, any)
	user           *types.UserDetails
}

type clientStreaming struct {
	authclient.AuthClient
	// *audioplayer.AudioPlayer
	*jobmanager.JobManager
	*syncer.NoSync

	*serverclient.ServerClient

	config Config
	log    *slog.Logger
	berr   bragerr.BragErrFactory

	eventCallbacks []func(types.ClientEvent, any)
	user           *types.UserDetails
}

type Client struct {
	clientFace
	*Identity
	playbackcontroller.PlaybackControllerFace

	*device.DeviceAgent

	activeRemoteDevice *uuid.UUID

	contextCallbacks  []func(context.Context, types.PlayContext)
	playbackCallbacks []func(context.Context, types.PlaybackState)

	sc  *serverclient.ServerClient
	log *slog.Logger
}

func (c *Client) SubscribeToClientEvents(handler sse.EventHandler) {
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

// FIXME ONLY FOR TEST
func (c *Client) LoginLocalUser(ctx context.Context, userID uuid.UUID) error {
	err := c.clientFace.LoginLocalUser(ctx, userID)
	if err != nil {
		return err
	}
	err = c.SubscribeDeviceEvents(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FIXME ONLY FOR TEST
func (c *Client) LoginToken(ctx context.Context, token string) (types.UserDetails, error) {
	user, err := c.clientFace.LoginToken(ctx, token)
	if err != nil {
		return types.UserDetails{}, err
	}
	err = c.SubscribeDeviceEvents(ctx)
	if err != nil {
		return types.UserDetails{}, err
	}
	return user, nil
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

func New(ctx context.Context, config Config, slogHandler slog.Handler) (ClientFace, error) {
	var cf clientFace
	var ar audioreader.AudioReader
	var err error

	sc := serverclient.New(config.ServerBaseURL, slogHandler)

	switch config.ClientType {
	case types.DeviceTypeStreaming:
		cf, err = newStreamingClient(ctx, config, &sc, slogHandler)
		if err != nil {
			return nil, err
		}
		ar = audioreader.NewServerReader(&sc, slogHandler)
	case types.DeviceTypeSync:
		cf, err = newSyncClient(ctx, config, &sc, slogHandler)
		if err != nil {
			return nil, err
		}
		ar = audioreader.NewLocalReader(config.MusicDirPath)
	default:
		return nil, errors.New("type not implemented")
	}

	id, err := NewIdentity("lucas test")
	if err != nil {
		return nil, err
	}

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		return nil, err
	}

	apCfg := audioplayer.Config{
		PlayerName:   config.PlayerName,
		MusicDirPath: config.MusicDirPath,
	}

	ap, err := audioplayer.New(apCfg, pa, ar, slogHandler)
	if err != nil {
		return nil, err
	}

	// FIXME: Needs to make another client structure. With authed starting the real clients
	da := device.NewAgent(slogHandler, &sc, uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")), types.DeviceBase{
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
	pc := playbackcontroller.New(ap, da, &sc, slogHandler)

	c := &Client{
		clientFace:             cf,
		Identity:               &id,
		PlaybackControllerFace: pc,
		sc:                     &sc,
		log:                    slog.New(slogHandler).With("service", "client"),
		DeviceAgent:            da,
		activeRemoteDevice:     nil,
		contextCallbacks:       []func(context.Context, types.PlayContext){},
		playbackCallbacks:      []func(context.Context, types.PlaybackState){},
	}

	// if err := da.SubscribeDeviceEvents(ctx); err != nil {
	// 	return nil, err
	// }
	ap.RegisterPlayCountCallback(c.updatePlayCount)
	ap.RegisterPlayContextCallback(c.handlePlayContextCallbacks)
	ap.RegisterPlaybackStateCallback(c.handlePlaybackStateCallbacks)

	da.SubscribeToEventCategory(c.handlePlayerEvents, "player")

	return c, nil
}
