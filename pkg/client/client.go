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
	"log/slog"

	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/audioreader"
	"github.com/bragemusic/core/pkg/authclient"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type Config struct {
	ConfigPath    string
	ImagePath     string
	MusicDirPath  string
	PlayerName    string
	ServerBaseURL string
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

	// LogoutServerUser logs out the currently authenticated server user
	// and invalidates any associated server session.
	LogoutServerUser(ctx context.Context) error

	// ServerStatus retrieves the current server API information and availability state.
	ServerStatus(ctx context.Context) (types.ServerApiInfo, error)
}

// AudioPlayerFace defines audio playback control and playback state
// observation functionality.
type AudioPlayerFace interface {
	// StartPlayerWithAlbum starts playback of an album beginning at the given track number.
	StartPlayerWithAlbum(ctx context.Context, albumID uuid.UUID, trackNumber int) error

	// StartPlayerWithPlaylist starts playback of a playlist beginning at the given track number.
	// Tracks are ordered according to sortBy and sortOrder.
	StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, sortBy database.SortBy, sortOrder database.SortOrder) error

	// AddTrackToQueue adds a track to the current playback queue.
	AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error

	// RegisterPlayContextChangeCallback registers a callback that is invoked
	// whenever the play context (album, playlist, queue, etc.) changes.
	RegisterPlayContextChangeCallback(f func(audioplayer.PlayContext))

	// RegisterPlayPauseCallback registers a callback that is invoked whenever
	// playback transitions between playing and paused states.
	RegisterPlayPauseCallback(f func(isPlaying bool))

	// RegisterProgressCallback registers a callback that is invoked periodically
	// with the current playback position in milliseconds.
	RegisterProgressCallback(f func(ms int64))

	// NextTrack skips to the next track in the current play context.
	NextTrack(ctx context.Context) (err error)

	// PlayContext returns the current playback context.
	PlayContext() audioplayer.PlayContext

	// PlayPause toggles playback between playing and paused states.
	PlayPause(ctx context.Context)

	// PreviousTrack skips to the previous track in the current play context.
	PreviousTrack(ctx context.Context) (err error)

	// SetRepeat sets the repeat mode for playback.
	SetRepeat(ctx context.Context, r audioplayer.RepeatType)

	// SetShuffle enables or disables shuffle mode for playback.
	SetShuffle(ctx context.Context, s bool)
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

	// ListUsers returns all users known to the system.
	ListUsers(ctx context.Context) ([]types.User, error)

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

// JobManagerFace defines background job execution and scheduling functionality.
type JobManagerFace interface {
	// StartScheduler starts the background job scheduler.
	// It typically runs until the context is canceled.
	StartScheduler(ctx context.Context)

	// RunJob executes a specific job type immediately.
	RunJob(ctx context.Context, jobType types.JobType) error
}

// ClientFace defines the high-level client interface that aggregates
// synchronization, authentication, playback, metadata, and job management
// functionality.
//
// It represents the main entry point for interacting with the system from
// a client perspective.
type ClientFace interface {
	// RegisterEventCallback registers a callback that is invoked whenever
	// a client-level evant is emitted. Messages may represent events,
	// warnings, or informational notifications originating from the client.
	RegisterEventCallback(f func(types.ClientEvent, any))

	SyncFace
	AuthFace
	AudioPlayerFace
	MetadataFace
	JobManagerFace
}

type ClientSync struct {
	*syncer.Syncer
	authclient.AuthClient
	*audioplayer.AudioPlayer
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

type ClientStreaming struct {
	authclient.AuthClient
	*audioplayer.AudioPlayer
	*jobmanager.JobManager
	*syncer.NoSync

	*serverclient.ServerClient

	config Config
	log    *slog.Logger
	berr   bragerr.BragErrFactory

	eventCallbacks []func(types.ClientEvent, any)
	user           *types.UserDetails
}

func NewSyncClient(ctx context.Context, config Config, slogHandler slog.Handler) (ClientFace, error) {
	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	mm := mediamanager.New(slogHandler, nil, nil, config.MusicDirPath, config.ImagePath)
	sy := syncer.New(&sc, nil, config.MusicDirPath, config.ImagePath, slogHandler)

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		return nil, err
	}

	ar := audioreader.NewLocalReader(config.MusicDirPath)

	apCfg := audioplayer.Config{
		PlayerName:   config.PlayerName,
		MusicDirPath: config.MusicDirPath,
	}

	ap, err := audioplayer.New(apCfg, pa, ar, slogHandler)
	if err != nil {
		return nil, err
	}

	jm := jobmanager.New(slogHandler)

	c := &ClientSync{
		config:       config,
		log:          slog.New(slogHandler).With("service", "client"),
		Syncer:       &sy,
		AuthClient:   authclient.New(&sc, slogHandler),
		AudioPlayer:  ap,
		MediaManager: &mm,
		JobManager:   &jm,
		sc:           &sc,
		berr:         bragerr.NewFactory("client"),
	}

	ap.RegisterPlayCountCallback(c.updatePlayCount)
	c.RegisterUserCallback(sy.SetUser)
	c.RegisterUserCallback(c.setUser)
	c.AuthClient.RegisterUpdateServerStatusCallback(c.updateServerStatusCallback)

	jm.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobAuthClientServerStatus,
		CronExpr: "*/10 * * * * *",
		Run:      c.AuthClient.UpdateServerStatus,
	})

	jm.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobSyncerDaemon,
		CronExpr: "*/10 * * * *",
		Run:      c.Syncer.Sync,
	})

	return c, nil
}

func NewStreamingClient(ctx context.Context, config Config, slogHandler slog.Handler) (ClientFace, error) {
	sc := serverclient.New(config.ServerBaseURL, slogHandler)

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		return nil, err
	}

	apCfg := audioplayer.Config{
		PlayerName:   config.PlayerName,
		MusicDirPath: config.MusicDirPath,
	}

	ar := audioreader.NewServerReader(&sc, slogHandler)

	ap, err := audioplayer.New(apCfg, pa, ar, slogHandler)
	if err != nil {
		return nil, err
	}

	jm := jobmanager.New(slogHandler)

	c := &ClientStreaming{
		config:       config,
		log:          slog.New(slogHandler).With("service", "client"),
		AuthClient:   authclient.New(&sc, slogHandler),
		AudioPlayer:  ap,
		JobManager:   &jm,
		ServerClient: &sc,
		NoSync:       &syncer.NoSync{},
		berr:         bragerr.NewFactory("client"),
	}

	ap.RegisterPlayCountCallback(c.updatePlayCount)
	c.RegisterUserCallback(c.setUser)
	// c.AuthClient.RegisterUpdateServerStatusCallback(c.updateServerStatusCallback)

	jm.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobAuthClientServerStatus,
		CronExpr: "*/10 * * * * *",
		Run:      c.AuthClient.UpdateServerStatus,
	})

	return c, nil
}
