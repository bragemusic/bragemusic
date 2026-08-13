package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/jobmanager"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/migrations"
	"github.com/bragemusic/core/pkg/musicbrainz"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/syncer"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (c *clientSync) RegisterEventCallback(f func(types.ClientEvent, any)) {
	c.eventCallbacks = append(c.eventCallbacks, f)
}

func (c *clientSync) emitEvent(event types.ClientEvent, payload any) {
	for _, f := range c.eventCallbacks {
		f(event, payload)
	}
}

func (c *clientSync) AddPlayCount(ctx context.Context, trackID uuid.UUID) error {
	return c.MediaManager.AddPlayCount(ctx, trackID, c.user.ID)
}

func (c *clientSync) Sync(ctx context.Context) error {
	if err := c.Syncer.Sync(ctx); err != nil {
		return err
	}

	c.emitEvent(types.ClientEventEntitiesUpdated, nil)
	return nil
}

func (c clientSync) GetArtistTopTracks(ctx context.Context, artistID uuid.UUID) ([]types.TrackDetailed, error) {
	tracks, err := c.ListTracksDetailedByArtist(ctx, artistID, c.user.ID, database.SortByPlayCount, database.SortDesc, utils.Ptr(10), false)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c clientSync) CreateArtist(ctx context.Context, artistData types.ArtistBase) error {
	if err := c.sc.CreateArtist(ctx, artistData); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) UpdateArtist(ctx context.Context, artistID uuid.UUID, artistData types.Artist) error {
	if err := c.sc.UpdateArtist(ctx, artistID, artistData); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) UploadArtistImage(ctx context.Context, artistID uuid.UUID, img serverclient.FileUpload) error {
	if err := c.sc.UploadArtistImage(ctx, artistID, img); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) ListTracksDetailedByAlbum(ctx context.Context, albumID uuid.UUID) ([]types.TrackDetailed, error) {
	return c.MediaManager.ListTracksDetailedByAlbum(ctx, albumID, c.user.ID)
}

func (c clientSync) FilterTracks(ctx context.Context, filter types.TrackFilter, page, limit int) (resp types.ListPaginationPayload[types.TrackDetailed], err error) {
	items, page, limit, totalPages, totalItems, err := c.ListTracksDetailed(ctx, filter, page, limit)
	if err != nil {
		return types.ListPaginationPayload[types.TrackDetailed]{}, err
	}

	return types.ListPaginationPayload[types.TrackDetailed]{
		Items:      items,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}, nil
}

func (c clientSync) UpdateAlbum(ctx context.Context, id uuid.UUID, album types.AlbumUpdate) error {
	if err := c.sc.UpdateAlbum(ctx, id, album); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) UploadAlbumImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.sc.UploadAlbumImage(ctx, id, img); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) RateTrack(ctx context.Context, trackID uuid.UUID, value int) error {
	if err := c.sc.RateTrack(ctx, trackID, value); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) UpdateTrack(ctx context.Context, id uuid.UUID, track types.TrackUpdate) error {
	if err := c.sc.UpdateTrack(ctx, id, track); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) AddPlaylist(ctx context.Context, playlist types.Playlist) error {
	if err := c.sc.AddPlaylist(ctx, playlist); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) AddSmartPlaylist(ctx context.Context, playlist types.PlaylistBase, filter types.TrackFilter) error {
	if err := c.sc.AddSmartPlaylist(ctx, playlist, filter); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) AddPlaylistTrack(ctx context.Context, playlistID, albumID, trackID uuid.UUID) error {
	if err := c.sc.AddPlaylistTrack(ctx, playlistID, albumID, trackID); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) CountPlaylists(ctx context.Context) (int, error) {
	return c.MediaManager.CountPlaylists(ctx, c.user.ID)
}

func (c clientSync) CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error) {
	return c.MediaManager.CountPlaylistTracks(ctx, playlistID, c.user.ID)
}

func (c clientSync) DeletePlaylist(ctx context.Context, id uuid.UUID) error {
	if err := c.sc.DeletePlaylist(ctx, id); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) DeletePlaylistTrack(ctx context.Context, id uuid.UUID) error {
	if err := c.sc.DeletePlaylistTrack(ctx, id); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) GetPlaylist(ctx context.Context, id uuid.UUID) (types.Playlist, error) {
	return c.MediaManager.GetPlaylist(ctx, id, c.user.ID)
}

func (c clientSync) ListPlaylists(ctx context.Context, includePublic bool, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.Playlist, error) {
	return c.MediaManager.ListPlaylists(ctx, c.user.ID, includePublic, sortBy, sortOrder)
}

func (c clientSync) ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID, sortBy database.SortBy, sortOrder database.SortOrder) ([]types.TrackDetailed, error) {
	return c.MediaManager.ListPlaylistTracks(ctx, playlistID, c.user.ID, sortBy, sortOrder)
}

func (c clientSync) UpdatePlaylist(ctx context.Context, id uuid.UUID, data types.Playlist) error {
	if err := c.sc.UpdatePlaylist(ctx, id, data); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) UploadPlaylistImage(ctx context.Context, id uuid.UUID, img serverclient.FileUpload) error {
	if err := c.sc.UploadPlaylistImage(ctx, id, img); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) ListEntityEvents(ctx context.Context) ([]types.EntityEvent, error) {
	return c.sc.ListEntityEvents(ctx)
}

func (c clientSync) CreateUser(ctx context.Context, email, username, password string, roles []types.UserRole) error {
	return c.sc.CreateUser(ctx, email, username, password, roles)
}

func (c clientSync) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return c.sc.DeleteUser(ctx, id)
}

func (c clientSync) ListUsers(ctx context.Context, includeMachineUsers bool) ([]types.UserDetails, error) {
	return c.sc.ListUsers(ctx, includeMachineUsers)
}

func (c clientSync) ListUserRoles(ctx context.Context) ([]types.UserRole, error) {
	return c.sc.ListUserRoles(ctx)
}

func (c clientSync) ListUserTokens(ctx context.Context) (tokens []types.TokenLimited, err error) {
	return c.sc.ListUserTokens(ctx)
}

func (c clientSync) DeleteUserToken(ctx context.Context, tokenID uuid.UUID) error {
	return c.sc.DeleteUserToken(ctx, tokenID)
}

func (c clientSync) CreateAPIToken(ctx context.Context, name string) (token string, err error) {
	return c.sc.CreateAPIToken(ctx, name)
}

func (c clientSync) UpdateUser(ctx context.Context, userID uuid.UUID, email, username string, password *string, roles []types.UserRole) error {
	return c.sc.UpdateUser(ctx, userID, email, username, password, roles)
}

func (c clientSync) UpdateProfile(ctx context.Context, data types.UpdateProfileReq) error {
	return c.sc.UpdateProfile(ctx, data)
}

func (c clientSync) UploadUserImage(ctx context.Context, r io.Reader, filename string) error {
	return c.sc.UploadUserImage(ctx, r, filename)
}

func (c clientSync) ImportAlbum(ctx context.Context, filename string, musicbrainzID *string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()
	return c.sc.ImportAlbum(ctx, r, filename, musicbrainzID)
}

func (c clientSync) ListImportItems(ctx context.Context, page, limit int) (types.ListPaginationPayload[types.Import], error) {
	return c.sc.ListImportItems(ctx, page, limit)
}

func (c clientSync) SearchMusicBrainz(ctx context.Context, artist, album string) (types.ListPayload[musicbrainz.SearchResults], error) {
	return c.sc.SearchMusicBrainz(ctx, artist, album)
}

func (c clientSync) GetTrackDetailed(ctx context.Context, trackID, albumID uuid.UUID) (track types.TrackDetailed, err error) {
	return c.MediaManager.GetTrackDetailed(ctx, trackID, albumID, c.user.ID)
}

func (c clientSync) LikeTrack(ctx context.Context, trackID uuid.UUID) error {
	if err := c.sc.LikeTrack(ctx, trackID); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) UnlikeTrack(ctx context.Context, trackID uuid.UUID) error {
	if err := c.sc.UnlikeTrack(ctx, trackID); err != nil {
		return err
	}
	return c.Sync(ctx)
}

func (c clientSync) ListLikedTracks(ctx context.Context) ([]types.TrackDetailed, error) {
	return c.MediaManager.ListLikedTracksDetailed(ctx, c.user.ID)
}

func (c clientSync) CountLikedTracks(ctx context.Context) (cnt int, err error) {
	tracks, err := c.MediaManager.ListLikedTracksDetailed(ctx, c.user.ID)
	if err != nil {
		return 0, err
	}

	return len(tracks), nil
}

func newSyncClient(ctx context.Context, config Config, jm *jobmanager.JobManager, sc *serverclient.ServerClient, user types.UserDetails, slogHandler slog.Handler) (clientFace, error) {
	dbPath := filepath.Join(config.ConfigPath, fmt.Sprintf("%s.db", user.ID.String()))
	if err := migrations.Migrate(ctx, dbPath, slogHandler); err != nil {
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

	mm := mediamanager.New(slogHandler, &db, nil, config.MusicDirPath, config.ImagePath)
	sy := syncer.New(sc, user, &db, config.MusicDirPath, config.ImagePath, slogHandler)

	c := &clientSync{
		config:       config,
		log:          slog.New(slogHandler).With("service", "client"),
		Syncer:       &sy,
		MediaManager: &mm,
		JobManager:   jm,
		user:         user,
		sc:           sc,
		berr:         bragerr.NewFactory("client"),
	}

	// FIXME: Live somewhere else
	// jm.RegisterJob(ctx, jobmanager.JobDefinition{
	// 	Type:     types.JobAuthClientServerStatus,
	// 	CronExpr: "*/10 * * * * *",
	// 	Run:      c.AuthClient.UpdateServerStatus,
	// })

	jm.RegisterJob(ctx, jobmanager.JobDefinition{
		Type:     types.JobSyncerDaemon,
		CronExpr: "*/10 * * * *",
		Run:      c.Syncer.Sync,
	})

	return c, nil
}
