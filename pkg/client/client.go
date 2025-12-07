package client

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/bragemusic/core/pkg/audiointerface"
	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/migrations"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
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
	*audioplayer.AudioPlayer
	sc      *serverclient.ServerClient
	mm      *mediamanager.MediaManager
	sy      *syncer.Syncer
	config  Config
	log     *slog.Logger
	dbClose func() error

	tracks []types.TrackEnhanced
}

func (c *Client) RegisterServerAvailabilityCallback(f func(bool)) {
	c.sy.RegisterServerAvailabilityCallback(f)
}

func (c *Client) RegisterSyncInProgressCallback(f func(bool)) {
	c.sy.RegisterSyncInProgressCallback(f)
}

func (c Client) Sync(ctx context.Context) error {
	return c.sy.Sync(ctx)
}

func (c *Client) Close() error {
	return c.dbClose()
}

func (c *Client) StartPlayer(ctx context.Context, albumID string, trackNumber int) error {
	tracks, err := c.mm.ListEnhancedTracksByAlbum(ctx, albumID)
	if err != nil {
		return err
	}

	c.tracks = tracks

	err = c.AudioPlayer.LoadAndStartTracks(ctx, c.tracks, trackNumber)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "albumID", albumID, "trackNumber", trackNumber)

	return nil
}

func (c Client) ListArtists(ctx context.Context) ([]types.Artist, error) {
	artists, err := c.mm.ListArtists(ctx)
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

func (c Client) ListAlbumsByArtist(ctx context.Context, artistID string) ([]types.Album, error) {
	albums, err := c.mm.ListAlbumsByArtist(ctx, artistID)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (c Client) GetAlbum(ctx context.Context, albumID string) (types.AlbumEnhanced, error) {
	album, err := c.mm.GetAlbumEnhanced(ctx, albumID)
	if err != nil {
		return types.AlbumEnhanced{}, err
	}
	return album, nil
}

func (c Client) ListTracksByAlbum(ctx context.Context, albumID string) ([]types.TrackEnhanced, error) {
	tracks, err := c.mm.ListEnhancedTracksByAlbum(ctx, albumID)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c *Client) StartSyncDaemon(ctx context.Context) {
	c.sy.StartDaemon(ctx)
}

// func (c Client) PlayPause(ctx context.Context) {
// 	c.ap.PlayPause(ctx)
// }

func NewSyncer(ctx context.Context, config Config, slogHandler slog.Handler) (c Client, err error) {
	dbPath := filepath.Join(config.ConfigPath, "data.db")
	if err = migrations.Migrate(ctx, dbPath, slogHandler); err != nil {
		return Client{}, err
	}

	dbSqlite, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return Client{}, err
	}

	db, err := database.New(dbSqlite)
	if err != nil {
		return Client{}, err
	}

	sc := serverclient.New(config.ServerBaseURL, slogHandler)
	mm := mediamanager.New(slogHandler, &db, config.MusicDirPath)
	sy := syncer.New(&sc, &db, config.MusicDirPath, config.ImagePath, slogHandler)

	pa, err := audiointerface.NewPortAudio(slogHandler)
	if err != nil {
		return Client{}, err
	}

	apCfg := audioplayer.Config{
		PlayerName:   config.PlayerName,
		MusicDirPath: config.MusicDirPath,
	}

	ap, err := audioplayer.New(apCfg, pa, slogHandler)
	if err != nil {
		return Client{}, err
	}

	return Client{
		sc:          &sc,
		mm:          &mm,
		sy:          &sy,
		AudioPlayer: ap,
		config:      config,
		log:         slog.New(slogHandler).With("service", "client"),
		dbClose:     dbSqlite.Close,
	}, nil
}
