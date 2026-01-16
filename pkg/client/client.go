package client

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/authclient"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/migrations"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	ConfigPath    string
	ImagePath     string
	MusicDirPath  string
	PlayerName    string
	ServerBaseURL string
}

type Client struct {
	authclient.AuthClient
	*audioplayer.AudioPlayer
	sc      *serverclient.ServerClient
	mm      *mediamanager.MediaManager
	sy      *syncer.Syncer
	config  Config
	log     *slog.Logger
	dbClose func() error

	// tracks []types.TrackEnhanced
}

func (c *Client) RegisterServerAvailabilityCallback(f func(server.Status)) {
	c.sy.RegisterServerAvailabilityCallback(f)
}

func (c *Client) RegisterSyncInProgressCallback(f func(bool)) {
	c.sy.RegisterSyncInProgressCallback(f)
}

func (c *Client) RegisterUserCallback(f func(*types.UserDetails)) {
	c.sy.RegisterUserCallback(f)
	c.AuthClient.RegisterUserCallback(f)
}

func (c Client) Sync(ctx context.Context) error {
	if err := c.sy.Sync(ctx); err != nil {
		return err
	}

	if err := c.sy.SyncItems(ctx); err != nil {
		return err
	}

	return nil
}

func (c Client) ServerStatus() server.Status {
	return c.sy.ServerStatus()
}

func (c *Client) Close() error {
	return c.dbClose()
}

func (c *Client) StartPlayerWithAlbum(ctx context.Context, albumID string, trackNumber int) error {
	albumUID, err := uuid.FromString(albumID)
	if err != nil {
		return err
	}

	tracks, err := c.mm.ListTracksDetailedByAlbum(ctx, albumUID)
	if err != nil {
		return err
	}

	pCtx := audioplayer.PlayContext{
		Type:            audioplayer.PlayContextAlbum,
		RefID:           albumID,
		Tracks:          tracks,
		Queue:           []types.TrackDetailed{},
		CurrentTrackIdx: trackNumber,
	}

	err = c.AudioPlayer.LoadAndStartTracks(ctx, pCtx)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "albumID", albumID, "trackNumber", trackNumber)

	return nil
}

func (c Client) ListArtists(ctx context.Context, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Artist, error) {
	artists, err := c.mm.ListArtists(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (c Client) GetArtist(ctx context.Context, artistID string) (types.Artist, error) {
	artist, err := c.mm.GetArtist(ctx, artistID)
	if err != nil {
		return types.Artist{}, err
	}
	return artist, nil
}

func (c Client) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	err := c.sc.UpdateArtist(ctx, artistID.String(), artistData)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) ListAlbumsByArtist(ctx context.Context, artistID string, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.AlbumDetailed, error) {
	albums, err := c.mm.ListAlbumsByArtist(ctx, artistID, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (c Client) GetAlbum(ctx context.Context, albumID string) (types.AlbumDetailed, error) {
	uid, err := uuid.FromString(albumID)
	if err != nil {
		return types.AlbumDetailed{}, err
	}

	album, err := c.mm.GetAlbumDetailed(ctx, uid)
	if err != nil {
		return types.AlbumDetailed{}, err
	}

	return album, nil
}

func (c Client) ListTracksByAlbum(ctx context.Context, albumID string) ([]types.TrackDetailed, error) {
	albumUID, err := uuid.FromString(albumID)
	if err != nil {
		return nil, err
	}

	tracks, err := c.mm.ListTracksDetailedByAlbum(ctx, albumUID)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c Client) UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error {
	err := c.sc.UpdateAlbum(ctx, id.String(), album)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) UploadAlbumImage(ctx context.Context, id string, img serverclient.ImageUpload) error {
	return c.sc.UploadAlbumImage(ctx, id, img)
}

func (c Client) UploadArtistImage(ctx context.Context, artistID string, img serverclient.ImageUpload) error {
	return c.sc.UploadArtistImage(ctx, artistID, img)
}

func (c Client) GetArtistTopTracks(ctx context.Context, artistID string) ([]types.TrackDetailed, error) {
	tracks, err := c.mm.ListTracksDetailedByArtist(ctx, artistID, database.SortByPlayCount, database.SortDesc, utils.Ptr(10), false)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c *Client) StartSyncDaemon(ctx context.Context, done func()) {
	c.sy.StartSyncDaemon(ctx, done)
}

func (c *Client) StartStatusDaemon(ctx context.Context, done func()) {
	c.sy.StartStatusDaemon(ctx, done)
}

func (c *Client) updatePlayCount(trackID string) {
	// FIXME: Do we need context here?
	ctx := context.TODO()
	// FIXME: Add proper UserID handling
	userID := "00000000-0000-0000-0000-000000000000"
	err := c.mm.AddPlayCount(ctx, trackID, userID)
	if err != nil {
		c.log.ErrorContext(ctx, "could not add play count", "error", err.Error())
		return
	}
	c.log.DebugContext(ctx, "added play count", "track_id", trackID)
}

// func (c Client) PlayPause(ctx context.Context) {
// 	c.ap.PlayPause(ctx)
// }

func NewSyncer(ctx context.Context, config Config, slogHandler slog.Handler) (c *Client, err error) {
	dbPath := filepath.Join(config.ConfigPath, "data.db")
	if err = migrations.Migrate(ctx, dbPath, slogHandler); err != nil {
		return nil, err
	}

	dbSqlite, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db, err := database.New(dbSqlite)
	if err != nil {
		return nil, err
	}

	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	mm := mediamanager.New(slogHandler, &db, nil, config.MusicDirPath, config.ImagePath)
	sy := syncer.New(&sc, &db, config.MusicDirPath, config.ImagePath, slogHandler)

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

	c = &Client{
		AuthClient:  authclient.New(&sc, slogHandler),
		sc:          &sc,
		mm:          &mm,
		sy:          &sy,
		AudioPlayer: ap,
		config:      config,
		log:         slog.New(slogHandler).With("service", "client"),
		dbClose:     dbSqlite.Close,
	}

	ap.RegisterPlayCountCallback(c.updatePlayCount)

	return c, nil
}
