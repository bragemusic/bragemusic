package client

import (
	"context"
	"log/slog"

	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/authclient"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type SyncFace interface {
	// Register a callback to know if a sync is running or not
	RegisterSyncInProgressCallback(f func(bool))
	// Returns if the client has the possibility to sync
	SupportsSync() bool
	// Sync server data to local client
	Sync(ctx context.Context) error
}

type AuthFace interface {
	RegisterUserCallback(f func(*types.UserDetails))
	LoginLocalUser(ctx context.Context, userID uuid.UUID) error
	LogoutLocalUser(ctx context.Context)
	LoginCachedServerUser(ctx context.Context, password string, longLivedToken bool) error
	GetUser() *types.UserDetails

	RegisterServerAvailabilityCallback(f func(server.ServerApiInfo))
	GetCachedUsers(ctx context.Context) (users []types.UserDetails, err error)
	Login(ctx context.Context, username, password string, longLivedToken bool) (types.UserDetails, error)
	LogoutServerUser(ctx context.Context) error
	ServerStatus(ctx context.Context) (server.ServerApiInfo, error)
	//
}

type AudioPlayerFace interface {
	StartPlayerWithAlbum(ctx context.Context, albumID uuid.UUID, trackNumber int) error
	StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, sortBy database.SortBy, sortOrder database.SortOrder) error
	AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error

	RegisterPlayContextChangeCallback(f func(audioplayer.PlayContext))
	RegisterPlayPauseCallback(f func(isPlaying bool))
	RegisterProgressCallback(f func(ms int64))

	NextTrack(ctx context.Context) (err error)
	PlayContext() audioplayer.PlayContext
	PlayPause(ctx context.Context)
	PreviousTrack(ctx context.Context) (err error)
	SetRepeat(ctx context.Context, r audioplayer.RepeatType)
	SetShuffle(ctx context.Context, s bool)
	//
}

type MetadataFace interface {
	CountArtists(ctx context.Context) (int, error)
	GetArtist(ctx context.Context, artistID uuid.UUID) (types.Artist, error)
	GetArtistTopTracks(ctx context.Context, artistID uuid.UUID) ([]types.TrackDetailed, error)
	ListArtists(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.ArtistDetailed, error)
	UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error          // sc
	UploadArtistImage(ctx context.Context, artistID uuid.UUID, img serverclient.FileUpload) error // sc

	CountAlbums(ctx context.Context) (int, error)
	GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (types.AlbumDetailed, error)
	ListAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error)
	ListAlbums(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error)
	UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error          // sc
	UploadAlbumImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error // sc

	CountTracks(ctx context.Context) (int, error)
	ListTracksDetailedByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.TrackDetailed, error)
	RateTrack(ctx context.Context, trackID uuid.UUID, value int) error            // sc
	UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error // sc

	AddPlaylist(ctx context.Context, playlist types.Playlist) error                     // sc
	AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error // sc
	CountPlaylists(ctx context.Context) (int, error)
	CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error)
	DeletePlaylist(ctx context.Context, id uuid.UUID) error      // sc
	DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error // sc
	GetPlaylist(ctx context.Context, id uuid.UUID) (types.Playlist, error)
	ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error)
	ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error)
	UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error              // sc
	UploadPlaylistImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error // sc

	ListEntityEvents(ctx context.Context) ([]types.EntityEvent, error) // sc

	ListUsers(ctx context.Context) ([]types.User, error) // sc

	SearchFull(ctx context.Context, searchTerm string) (si []types.SearchItem, err error)

	ImportAlbum(ctx context.Context, filename string, musicbrainzID *string) error // sc, wrapper

	AddPlayCount(ctx context.Context, trackID, userID uuid.UUID) error
}

type JobManagerFace interface {
	StartScheduler(ctx context.Context)
	RunJob(ctx context.Context, jobType types.JobType) error
}

kolla pa en client o authclient. Chad

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

	user *types.UserDetails
}

func NewSyncClient(ctx context.Context, config Config, slogHandler slog.Handler) (ClientFace, error) {
	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	mm := mediamanager.New(slogHandler, nil, nil, config.MusicDirPath, config.ImagePath)
	sy := syncer.New(&sc, nil, config.MusicDirPath, config.ImagePath, slogHandler)

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		return nil, err
	}

	apCfg := audioplayer.Config{
		PlayerName:   config.PlayerName,
		MusicDirPath: config.MusicDirPath,
	}

	ap, err := audioplayer.New(apCfg, pa, slogHandler)
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

type ClientFace interface {
	RegisterMsgCallback(f func(types.ClientMessage))
	SyncFace
	AuthFace
	AudioPlayerFace
	MetadataFace
	JobManagerFace
	//
}

type Config struct {
	ConfigPath    string
	ImagePath     string
	MusicDirPath  string
	PlayerName    string
	ServerBaseURL string
}

type Client struct {
	// authclient.AuthClient
	// *audioplayer.AudioPlayer
	*jobmanager.JobManager
	sc *serverclient.ServerClient
	mm *mediamanager.MediaManager
	// sy      *syncer.Syncer
	config  Config
	log     *slog.Logger
	dbClose func() error
	// FIXME: This should not be here later. Just now for the search job. when the jobs is its own service it should hold the db
	db database.DatabaseFace

	// tracks []types.TrackEnhanced
}

// func (c *Client) setDatabase(ctx context.Context, dbPath string) error {
// 	if c.dbClose != nil {
// 		c.log.InfoContext(ctx, "closing database")
// 		if err := c.dbClose(); err != nil {
// 			return err
// 		}
// 	}

// 	if err := migrations.Migrate(ctx, dbPath, c.log.Handler()); err != nil {
// 		return err
// 	}

// 	c.log.InfoContext(ctx, "opening database", "path", dbPath)
// 	dbSqlite, err := sqlx.Open("sqlite3", dbPath)
// 	if err != nil {
// 		return err
// 	}

// 	db, err := database.New(dbSqlite)
// 	if err != nil {
// 		return err
// 	}

// 	c.mm.SetDatabase(&db)
// 	c.sy.SetDatabase(&db)

// 	c.db = db
// 	c.dbClose = dbSqlite.Close

// 	return nil
// }

// func (c *Client) LoginLocalUser(ctx context.Context, userID uuid.UUID) error {
// 	c.log.InfoContext(ctx, "logging in local user", "id", userID.String())

// 	user, err := c.AuthClient.LoginLocalUser(ctx, userID, false)
// 	if err != nil {
// 		return err
// 	}

// 	dbPath := filepath.Join(c.config.ConfigPath, fmt.Sprintf("%s.db", userID.String()))

// 	if err := c.setDatabase(ctx, dbPath); err != nil {
// 		return err
// 	}

// 	c.AuthClient.UserCallback(&user)
// 	c.ServerStatus(ctx)

// 	return nil
// }

// func (c Client) PlayPause(ctx context.Context) {
// 	c.ap.PlayPause(ctx)
// }

// func NewSyncer(ctx context.Context, config Config, slogHandler slog.Handler) (c *Client, err error) {
// 	sc := serverclient.New(config.ServerBaseURL, slogHandler)
// 	mm := mediamanager.New(slogHandler, nil, nil, config.MusicDirPath, config.ImagePath)
// 	// sy := syncer.New(&sc, nil, config.MusicDirPath, config.ImagePath, slogHandler)

// 	pa, err := audiointerface.NewPortAudio(slogHandler)
// 	if err != nil {
// 		return nil, err
// 	}

// 	apCfg := audioplayer.Config{
// 		PlayerName:   config.PlayerName,
// 		MusicDirPath: config.MusicDirPath,
// 	}

// 	ap, err := audioplayer.New(apCfg, pa, slogHandler)
// 	if err != nil {
// 		return nil, err
// 	}

// 	jm := jobmanager.New(slogHandler)

// 	c = &Client{
// 		// AuthClient: authclient.New(&sc, slogHandler),
// 		sc: &sc,
// 		mm: &mm,
// 		// sy:          &sy,
// 		JobManager: &jm,
// 		// AudioPlayer: ap,
// 		config: config,
// 		log:    slog.New(slogHandler).With("service", "client"),
// 	}

// 	ap.RegisterPlayCountCallback(c.updatePlayCount)
// 	// c.RegisterUserCallback(sy.SetUser)
// 	// c.AuthClient.RegisterUpdateServerStatusCallback(c.updateServerStatusCallback)

// 	// jm.RegisterJob(ctx, jobmanager.JobDefinition{
// 	// 	Type:     types.JobAuthClientServerStatus,
// 	// 	CronExpr: "*/10 * * * * *",
// 	// 	Run:      c.AuthClient.UpdateServerStatus,
// 	// })

// 	// jm.RegisterJob(ctx, jobmanager.JobDefinition{
// 	// 	Type:     types.JobSyncerDaemon,
// 	// 	CronExpr: "*/10 * * * *",
// 	// 	Run:      c.sy.Daemon,
// 	// })

// 	return c, nil
// }
