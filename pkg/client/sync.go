package client

import (
	"context"
	"errors"
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

func (c *ClientSync) RegisterEventCallback(f func(types.ClientEvent, any)) {
	c.eventCallbacks = append(c.eventCallbacks, f)
}

func (c *ClientSync) emitEvent(event types.ClientEvent, payload any) {
	for _, f := range c.eventCallbacks {
		f(event, payload)
	}
}

func (c *ClientSync) setUser(user *types.UserDetails) {
	c.user = user
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

func (c *ClientSync) AddPlayCount(ctx context.Context, trackID uuid.UUID) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not get top tracks"))
	}

	return c.MediaManager.AddPlayCount(ctx, trackID, c.user.ID)
}

func (c *ClientSync) Sync(ctx context.Context) error {
	if err := c.Syncer.Sync(ctx); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
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

func (c *ClientSync) LoginCachedServerUser(ctx context.Context, password string, longLivedToken bool) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not login cached user"))
	}

	return c.AuthClient.LoginCachedServerUser(ctx, *c.user, password, longLivedToken)
}

func (c *ClientSync) LogoutLocalUser(ctx context.Context) {
	c.AuthClient.UserCallback(nil)
}

func (c *ClientSync) LogoutServerUser(ctx context.Context) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not logout server user"))
	}
	return c.AuthClient.LogoutServerUser(ctx, c.user.ID)
}

func (c ClientSync) GetUser() *types.UserDetails {
	return c.user
}

func (c ClientSync) GetArtistTopTracks(ctx context.Context, artistID uuid.UUID) ([]types.TrackDetailed, error) {
	if c.user == nil {
		return nil, c.berr.NoUserInContext(errors.New("could not get top tracks"))
	}

	tracks, err := c.ListTracksDetailedByArtist(ctx, artistID, c.user.ID, database.SortByPlayCount, database.SortDesc, utils.Ptr(10), false)
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

func (c ClientSync) ListTracksDetailedByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.TrackDetailed, error) {
	if c.user == nil {
		return nil, c.berr.NoUserInContext(errors.New("could not list tracks"))
	}
	return c.MediaManager.ListTracksDetailedByAlbum(ctx, albumID, c.user.ID)
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

func (c ClientSync) CountPlaylists(ctx context.Context) (int, error) {
	if c.user == nil {
		return 0, c.berr.NoUserInContext(errors.New("could count playlists"))
	}
	return c.MediaManager.CountPlaylists(ctx, c.user.ID)
}

func (c ClientSync) CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error) {
	if c.user == nil {
		return 0, c.berr.NoUserInContext(errors.New("could count playlist tracks"))
	}
	return c.MediaManager.CountPlaylistTracks(ctx, playlistID, c.user.ID)
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

func (c ClientSync) GetPlaylist(ctx context.Context, id uuid.UUID) (types.Playlist, error) {
	if c.user == nil {
		return types.Playlist{}, c.berr.NoUserInContext(errors.New("could get playlist"))
	}
	return c.MediaManager.GetPlaylist(ctx, id, c.user.ID)
}

func (c ClientSync) ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error) {
	if c.user == nil {
		return nil, c.berr.NoUserInContext(errors.New("could list playlists"))
	}
	return c.MediaManager.ListPlaylists(ctx, c.user.ID, includePublic, sortBy, sortOrder)
}

func (c ClientSync) ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error) {
	if c.user == nil {
		return nil, c.berr.NoUserInContext(errors.New("could list playlist tracks"))
	}
	return c.MediaManager.ListPlaylistTracks(ctx, playlistID, c.user.ID, sortBy, sortOrder)
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

func (c *ClientSync) StartPlayerWithAlbum(ctx context.Context, albumID uuid.UUID, trackNumber int) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not start player"))
	}

	tracks, err := c.MediaManager.ListTracksDetailedByAlbum(ctx, albumID, c.user.ID)
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

func (c *ClientSync) StartPlayerWithLikedTracks(ctx context.Context, trackNumber int) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not start player"))
	}

	tracks, err := c.ListLikedTracks(ctx)
	if err != nil {
		return err
	}

	pCtx := audioplayer.PlayContext{
		Type:            audioplayer.PlayContextLikedTracks,
		RefID:           uuid.Nil,
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

	c.log.InfoContext(ctx, "started player", "type", "liked tracks", "trackNumber", trackNumber)

	return nil
}

func (c *ClientSync) StartPlayerWithPlaylist(ctx context.Context, playlistID uuid.UUID, trackNumber int, sortBy database.SortBy, sortOrder database.SortOrder) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not start player"))
	}

	tracks, err := c.MediaManager.ListPlaylistTracks(ctx, playlistID, c.user.ID, sortBy, sortOrder)
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

func (c *ClientSync) AddTrackToQueue(ctx context.Context, trackID, albumID uuid.UUID) error {
	if c.user == nil {
		return c.berr.NoUserInContext(errors.New("could not start player"))
	}

	track, err := c.MediaManager.GetTrackDetailed(ctx, trackID, albumID, c.user.ID)
	if err != nil {
		return err
	}

	c.AudioPlayer.AddTrackToQueue(ctx, track)
	return nil
}

func (c ClientSync) LikeTrack(ctx context.Context, trackID uuid.UUID) error {
	if err := c.sc.LikeTrack(ctx, trackID); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) UnlikeTrack(ctx context.Context, trackID uuid.UUID) error {
	if err := c.sc.UnlikeTrack(ctx, trackID); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c ClientSync) ListLikedTracks(ctx context.Context) ([]types.TrackDetailed, error) {
	if c.user == nil {
		return nil, c.berr.NoUserInContext(errors.New("could not list liked tracks"))
	}
	return c.MediaManager.ListLikedTracksDetailed(ctx, c.user.ID)
}

func (c ClientSync) CountLikedTracks(ctx context.Context) (cnt int, err error) {
	if c.user == nil {
		return 0, c.berr.NoUserInContext(errors.New("could not count liked tracks"))
	}

	tracks, err := c.MediaManager.ListLikedTracksDetailed(ctx, c.user.ID)
	if err != nil {
		return 0, err
	}

	return len(tracks), nil
}
