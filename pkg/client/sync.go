package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/audioplayer"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/migrations"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (c *ClientSync) RegisterMsgCallback(f func(types.ClientMessage)) {
}

func (c *ClientSync) updatePlayCount(trackID uuid.UUID) {
	// FIXME: Do we need context here?
	ctx := context.TODO()

	if c.user == nil {
		c.log.ErrorContext(ctx, "could not update play count, no user in context")
		return
	}

	err := c.MediaManager.AddPlayCount(ctx, trackID, c.user.ID)
	if err != nil {
		c.log.ErrorContext(ctx, "could not add play count", "error", err.Error())
		return
	}
	c.log.DebugContext(ctx, "added play count", "track_id", trackID)
}

func (c *ClientSync) updateServerStatusCallback(ctx context.Context) {
	c.ServerStatus(ctx)
}

func (c *ClientSync) setDatabase(ctx context.Context, dbPath string) error {
	if c.dbClose != nil {
		c.log.InfoContext(ctx, "closing database")
		if err := c.dbClose(); err != nil {
			return err
		}
	}
	if err := migrations.Migrate(ctx, dbPath, c.log.Handler()); err != nil {
		return err
	}

	c.log.InfoContext(ctx, "opening database", "path", dbPath)
	dbSqlite, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	db, err := database.New(dbSqlite)
	if err != nil {
		return err
	}

	c.MediaManager.SetDatabase(&db)
	c.Syncer.SetDatabase(&db)

	c.dbClose = dbSqlite.Close

	return nil
}

func (c *ClientSync) LoginLocalUser(ctx context.Context, userID uuid.UUID) error {
	c.log.InfoContext(ctx, "logging in local user", "id", userID.String())

	user, err := c.AuthClient.LoginLocalUser(ctx, userID, false)
	if err != nil {
		return err
	}

	dbPath := filepath.Join(c.config.ConfigPath, fmt.Sprintf("%s.db", userID.String()))

	if err := c.setDatabase(ctx, dbPath); err != nil {
		return err
	}

	c.AuthClient.UserCallback(&user)
	c.ServerStatus(ctx)

	return nil
}

func (c ClientSync) GetArtistTopTracks(ctx context.Context, artistID, userID uuid.UUID) ([]types.TrackDetailed, error) {
	tracks, err := c.ListTracksDetailedByArtist(ctx, artistID, userID, database.SortByPlayCount, database.SortDesc, utils.Ptr(10), false)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c ClientSync) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	if err := c.sc.UpdateArtist(ctx, artistID, artistData); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UploadArtistImage(ctx context.Context, artistID uuid.UUID, img serverclient.FileUpload) error {
	if err := c.sc.UploadArtistImage(ctx, artistID, img); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error {
	if err := c.sc.UpdateAlbum(ctx, id, album); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UploadAlbumImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.sc.UploadAlbumImage(ctx, id, img); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) RateTrack(ctx context.Context, trackID uuid.UUID, value int) error {
	if err := c.sc.RateTrack(ctx, trackID, value); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error {
	if err := c.sc.UpdateTrack(ctx, id, track); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) AddPlaylist(ctx context.Context, playlist types.Playlist) error {
	if err := c.sc.AddPlaylist(ctx, playlist); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error {
	if err := c.sc.AddPlaylistTrack(ctx, playlistID, albumID, trackID); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) DeletePlaylist(ctx context.Context, id uuid.UUID) error {
	if err := c.sc.DeletePlaylist(ctx, id); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error {
	if err := c.sc.DeletePlaylistTrack(ctx, id); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error {
	if err := c.sc.UpdatePlaylist(ctx, id, data); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UploadPlaylistImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.sc.UploadPlaylistImage(ctx, id, img); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) ListEntityEvents(ctx context.Context) ([]types.EntityEvent, error) {
	return c.sc.ListEntityEvents(ctx)
}

func (c ClientSync) ListUsers(ctx context.Context) ([]types.User, error) {
	return c.sc.ListUsers(ctx)
}

func (c ClientSync) ImportAlbum(ctx context.Context, filename string, musicbrainzID *string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()
	return c.sc.ImportAlbum(ctx, r, filename, musicbrainzID)
}

func (c *ClientSync) StartPlayerWithAlbum(ctx context.Context, userID, albumID uuid.UUID, trackNumber int) error {
	tracks, err := c.MediaManager.ListTracksDetailedByAlbum(ctx, albumID, userID)
	if err != nil {
		return err
	}

	pCtx := audioplayer.PlayContext{
		Type:            audioplayer.PlayContextAlbum,
		RefID:           albumID,
		Tracks:          tracks,
		Queue:           []types.TrackDetailed{},
		CurrentTrackIdx: trackNumber,
		Shuffle:         c.PlayContext().Shuffle,
		Repeat:          c.PlayContext().Repeat,
	}

	err = c.AudioPlayer.LoadAndStartTracks(ctx, pCtx)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "albumID", albumID.String(), "trackNumber", trackNumber)

	return nil
}

func (c *ClientSync) StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, userID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) error {
	tracks, err := c.MediaManager.ListPlaylistTracks(ctx, playlistID, userID, sortBy, sortOrder)
	if err != nil {
		return err
	}

	pCtx := audioplayer.PlayContext{
		Type:            audioplayer.PlayContextPlaylist,
		RefID:           playlistID,
		Tracks:          tracks,
		Queue:           []types.TrackDetailed{},
		CurrentTrackIdx: trackNumber,
		Shuffle:         c.PlayContext().Shuffle,
		Repeat:          c.PlayContext().Repeat,
	}

	err = c.AudioPlayer.LoadAndStartTracks(ctx, pCtx)
	if err != nil {
		return err
	}

	c.log.InfoContext(ctx, "started player", "playlistID", playlistID.String(), "trackNumber", trackNumber)

	return nil
}

func (c *ClientSync) AddTrackToQueue(ctx context.Context, trackID, albumID, userID uuid.UUID) error {
	track, err := c.MediaManager.GetTrackDetailed(ctx, trackID, albumID, userID)
	if err != nil {
		return err
	}

	c.AudioPlayer.AddTrackToQueue(ctx, track)
	return nil
}
