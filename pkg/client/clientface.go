package client

import (
	"context"
	"io"

	"github.com/bragemusic/bragemusic/pkg/database"
	"github.com/bragemusic/bragemusic/pkg/musicbrainz"
	"github.com/bragemusic/bragemusic/pkg/serverclient"
	"github.com/bragemusic/bragemusic/pkg/sse"
	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
)

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
	GetConnectedDevice() *uuid.UUID
}

// MetadataFace defines access to and modification of music library metadata,
// including artists, albums, tracks, playlists, users, and search functionality.
type MetadataFace interface {
	// CountArtists returns the total number of artists.
	CountArtists(ctx context.Context) (int, error)

	// CreateArtist creates a new Artist object
	CreateArtist(ctx context.Context, artistData types.ArtistBase) error

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

	// ListFeaturedAlbumsByArtist returns albums where a given artist is featured.
	ListFeaturedAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) (albums []types.AlbumDetailed, err error)

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

	// FilterTracks returns detailed track metadata of tracks matching the filter
	FilterTracks(ctx context.Context, filter types.TrackFilter, page, limit int) (types.ListPaginationPayload[types.TrackDetailed], error)

	// RateTrack sets the rating value for a specific track.
	RateTrack(ctx context.Context, trackID uuid.UUID, value int) error

	// UpdateTrack updates metadata for a track.
	UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error

	// AddPlaylist creates a new playlist.
	AddPlaylist(ctx context.Context, playlist types.Playlist) error

	// AddSmartPlaylist creates a new smart playlist.
	AddSmartPlaylist(ctx context.Context, playlist types.PlaylistBase, filter types.TrackFilter) error

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

	// ListImportItems returns all import items, queued, progressing, finished and errored.
	ListImportItems(ctx context.Context, page, limit int) (types.ListPaginationPayload[types.Import], error)

	// SearchMusicBrainz searches after a release id on musicbrainz using artist and or album as input
	SearchMusicBrainz(ctx context.Context, artist, album string) (types.ListPayload[musicbrainz.SearchResults], error)

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

	// CreateAPIToken creates an API token for the logged in user.
	CreateAPIToken(ctx context.Context, name string) (token string, err error)

	// DeleteUser removes the selected user from the server. This action cannot be reversed.
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// DeleteUserToken removes a specified token owned by the logged in user.
	DeleteUserToken(ctx context.Context, tokenID uuid.UUID) error

	// ListUsers returns all users known to the system.
	ListUsers(ctx context.Context, includeMachineUsers bool) ([]types.UserDetails, error)

	// ListUserRoles returns a list of available user roles
	ListUserRoles(ctx context.Context) ([]types.UserRole, error)

	// ListUserTokens returns a list of tokens belonging to the logged in user.
	ListUserTokens(ctx context.Context) (tokens []types.TokenLimited, err error)

	// UpdateUser updates the selected user's information. If password is not nil, it will be changed.
	UpdateUser(ctx context.Context, userID uuid.UUID, email, username string, password *string, roles []types.UserRole) error

	// UpdateProfile updates the currently logged in users profile information.
	UpdateProfile(ctx context.Context, data types.UpdateProfileReq) error

	// UploadUserImage uploads a profile picture for the logged in user.
	UploadUserImage(ctx context.Context, r io.Reader, filename string) error
	//
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
	// AuthFace
	MetadataFace
	JobManagerFace
	AccessFace
}

type ClientFace interface {
	SubscribeToClientEvents(handler sse.EventHandler)
	RegisterServerAvailabilityCallback(f func(types.ServerApiInfo))
	Close(context.Context) error
	ServerStatus(ctx context.Context) (types.ServerApiInfo, error)
	clientFace
	IdentityFace
	AudioPlayerFace
	DeviceManagerFace

	// hit ska åtminstone auth och audioplayer flyttas. Det är delat mellan sync och streaming. Typ iaf. Får lösa så att auth är det.
	// 	Sen ska clientface generarea en clientSync/Stream när det loggas in en användare, så slipper vi pekar på user o en massa skit
}
