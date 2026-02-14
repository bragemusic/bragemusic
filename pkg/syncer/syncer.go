package syncer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

const expectedServerApplication = "brage-server"

type Syncer struct {
	sc                       *serverclient.ServerClient
	db                       database.DatabaseFace
	log                      *slog.Logger
	musicDir                 string
	imgDir                   string
	serverAvailable          bool
	serverStatus             server.Status
	syncInProgress           bool
	user                     *types.UserDetails
	serverAvailableCallbacks []func(server.Status)
	syncInProgressCallbacks  []func(bool)
	userCallbacks            []func(*types.UserDetails)
}

func (s *Syncer) RegisterServerAvailabilityCallback(f func(server.Status)) {
	s.serverAvailableCallbacks = append(s.serverAvailableCallbacks, f)
}

func (s *Syncer) RegisterSyncInProgressCallback(f func(bool)) {
	s.syncInProgressCallbacks = append(s.syncInProgressCallbacks, f)
}

func (s *Syncer) RegisterUserCallback(f func(*types.UserDetails)) {
	s.userCallbacks = append(s.userCallbacks, f)
}

func (s *Syncer) StartStatusDaemon(ctx context.Context, done func()) {
	go func() {
		defer done()

		tickerStatus := time.NewTicker(10 * time.Second)
		defer tickerStatus.Stop()

		for {
			select {
			case <-tickerStatus.C:
				s.log.DebugContext(ctx, "updating server availability")
				err := s.updateServerAvailability(ctx)
				if err != nil {
					s.log.WarnContext(ctx, "server unreachable", "error", err.Error())
					s.serverAvailable = false
				} else {
					s.serverAvailable = true
					if s.user == nil {
						user, err := s.sc.GetUser(ctx)
						if err != nil {
							s.log.ErrorContext(ctx, err.Error())
						} else {
							s.user = &user
							for _, f := range s.userCallbacks {
								f(s.user)
							}
						}
					}
				}

			case <-ctx.Done():
				s.log.InfoContext(ctx, "terminating status check")
				return
			}
		}
	}()
}

func (s *Syncer) StartSyncDaemon(ctx context.Context, done func()) {
	go func() {
		defer done()

		tickerSync := time.NewTicker(15 * time.Minute)
		defer tickerSync.Stop()

		for {
			select {
			case <-tickerSync.C:
				if s.serverAvailable {
					s.log.InfoContext(ctx, "starting periodic sync")
					err := s.Sync(ctx)
					if err != nil {
						s.log.ErrorContext(ctx, "periodic sync finished with errors", "error", err.Error())
					}
					err = s.SyncItems(ctx)
					if err != nil {
						s.log.ErrorContext(ctx, "periodic sync finished with errors", "error", err.Error())
					}
				} else {
					s.log.DebugContext(ctx, "periodic sync skipped. Server not available")
				}

			case <-ctx.Done():
				s.log.InfoContext(ctx, "terminating periodic sync")
				return
			}
		}
	}()
}

func (s *Syncer) updateServerAvailability(ctx context.Context) error {
	h, err := s.sc.CheckStatus(ctx)
	if err != nil {
		h = server.Status{
			Status: server.HealthzUnavailable,
		}
	}

	s.serverStatus = h

	for _, f := range s.serverAvailableCallbacks {
		f(h)
	}
	if err != nil {
		return err
	}

	if h.Application != expectedServerApplication {
		return fmt.Errorf("expected server application name differs. '%s' != expected '%s'", h.Application, expectedServerApplication)
	}

	if h.Status == server.HealthzUnavailable {
		return fmt.Errorf("server status is not running. '%s'", h.Status)
	}

	return nil
}

func (s *Syncer) ServerStatus() server.Status {
	return s.serverStatus
}

func (s *Syncer) Sync(ctx context.Context) error {
	if !s.serverAvailable {
		return errors.New("server is not available")
	}

	s.syncInProgress = true
	for _, f := range s.syncInProgressCallbacks {
		f(true)
	}

	defer func() {
		s.syncInProgress = false
		for _, f := range s.syncInProgressCallbacks {
			f(false)
		}
	}()

	lastSync, err := s.db.GetLastSync(ctx)
	if err != nil {
		lastSync.SyncedAt = time.Unix(0, 0)
		s.log.InfoContext(ctx, "starting sync, no previous sync found")
	} else {
		s.log.InfoContext(ctx, fmt.Sprintf("starting sync, syncing items since %s", lastSync.SyncedAt.String()))
	}

	syncState, err := s.sc.GetSyncState(ctx, lastSync.SyncedAt)
	if err != nil {
		return err
	}

	dbSyncState := types.DBSyncState{}
	dbSyncState.SyncedAt = syncState.Time

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	dbSyncState.ArtistsCreated, dbSyncState.ArtistsUpdated, err = s.syncArtists(ctx, tx, syncState.CreatedOrUpdated.Artists)
	if err != nil {
		return err
	}

	dbSyncState.AlbumsCreated, dbSyncState.AlbumsUpdated, err = s.syncAlbums(ctx, tx, syncState.CreatedOrUpdated.Albums)
	if err != nil {
		return err
	}

	dbSyncState.TracksCreated, dbSyncState.TracksUpdated, err = s.syncTracks(ctx, tx, syncState.CreatedOrUpdated.Tracks)
	if err != nil {
		return err
	}

	_, _, err = s.syncAlbumArtists(ctx, tx, syncState.CreatedOrUpdated.AlbumArtists, syncState.Deleted.AlbumArtists)
	if err != nil {
		return err
	}

	_, _, err = s.syncAlbumTracks(ctx, tx, syncState.CreatedOrUpdated.AlbumTracks)
	if err != nil {
		return err
	}

	_, _, err = s.syncPlaylists(ctx, tx, syncState.CreatedOrUpdated.Playlists, syncState.Deleted.Playlists)
	if err != nil {
		return err
	}

	_, _, err = s.syncPlaylistTracks(ctx, tx, syncState.CreatedOrUpdated.PlaylistTracks, syncState.Deleted.PlaylistTracks)
	if err != nil {
		return err
	}

	if err = s.syncMediaFiles(ctx, tx, syncState.CreatedOrUpdated.MediaFiles); err != nil {
		return err
	}

	if err = s.syncEntityEvents(ctx, tx, syncState.New); err != nil {
		return err
	}

	if err := s.syncPlayHistory(ctx, tx, lastSync.SyncedAt); err != nil {
		return err
	}

	if _, err := tx.AddSync(ctx, dbSyncState); err != nil {
		return err
	}

	totalCreations := dbSyncState.ArtistsCreated + dbSyncState.AlbumsCreated + dbSyncState.TracksCreated
	totalUpdates := dbSyncState.ArtistsUpdated + dbSyncState.AlbumsUpdated + dbSyncState.TracksUpdated

	if totalCreations == 0 && totalUpdates == 0 {
		s.log.InfoContext(ctx, "sync finished. no entries added or updated")
	} else {
		s.log.InfoContext(ctx, fmt.Sprintf("sync finished. %d entries added and %d updated", totalCreations, totalUpdates))
	}

	return tx.Commit()
}

func (s *Syncer) syncEntityEvents(ctx context.Context, tx database.DatabaseFace, events []types.EntityEvent) error {
	for _, e := range events {
		var f func(ctx context.Context, tx database.DatabaseFace, event types.EntityEvent) error

		switch e.EntityType {
		case types.EntityRating:
			f = s.syncRating
		default:
			return fmt.Errorf("unsupported entity type '%s'", e.EntityType)
		}

		if err := f(ctx, tx, e); err != nil {
			return err
		}

	}
	return nil
}

func (s *Syncer) syncRating(ctx context.Context, tx database.DatabaseFace, event types.EntityEvent) error {
	switch event.Type {
	case types.EntityEventCreate:
		r, err := s.sc.GetRating(ctx, event.ItemID)
		if err != nil {
			return err
		}
		if _, err := tx.AddRating(ctx, r); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		r, err := s.sc.GetRating(ctx, event.ItemID)
		if err != nil {
			return err
		}
		if err := tx.UpdateRating(ctx, r.ID, r.Rating); err != nil {
			return err
		}
	case types.EntityEventDelete:
		return errors.New("'delete' not supported for ratings")
	}
	return nil
}

func (s *Syncer) SyncItems(ctx context.Context) error {
	if !s.serverAvailable {
		return errors.New("server is not available")
	}

	s.syncInProgress = true
	for _, f := range s.syncInProgressCallbacks {
		f(true)
	}

	defer func() {
		s.syncInProgress = false
		for _, f := range s.syncInProgressCallbacks {
			f(false)
		}
	}()

	for {
		tx, err := s.db.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		si, err := tx.GetUnsyncedItem(ctx)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				break
			}
			return err
		}

		if si.Type != types.SiTypeMediaFile {
			return errors.New("only media files can be asyncronically synced right now")
		}

		mf, err := s.sc.GetMediaFile(ctx, si.ItemID)
		if err != nil {
			return err
		}

		exists, err := tx.MediaFileExists(ctx, mf.ID)
		if err != nil {
			return err
		}

		if exists {
			if err = tx.UpdateMediaFile(ctx, mf); err != nil {
				return err
			}
		} else {
			if _, err = tx.AddMediaFile(ctx, mf); err != nil {
				return err
			}
		}

		mfPath := filepath.Join(s.musicDir, mf.Filename())

		if err = os.MkdirAll(filepath.Dir(mfPath), os.ModePerm); err != nil {
			return err
		}

		dst, err := os.Create(mfPath)
		if err != nil {
			return err
		}

		s.log.InfoContext(ctx, fmt.Sprintf("downloading media file '%s' to '%s'", mf.ID.String(), mfPath))

		if err = s.sc.DownloadMediaFile(ctx, mf.ID, dst); err != nil {
			dst.Close()
			return err
		}

		if err = dst.Close(); err != nil {
			return err
		}

		if err = tx.SetSyncItemState(ctx, si.ID, types.SiStateFinished); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (s Syncer) syncArtists(ctx context.Context, tx database.DatabaseFace, artistIDs []string) (created, updated int, err error) {
	for _, aID := range artistIDs {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing artist '%s'", aID))
		exists, err := tx.ArtistExists(ctx, aID)
		if err != nil {
			return 0, 0, err
		}

		serverArtist, err := s.sc.GetArtist(ctx, aID)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			if err = tx.UpdateArtist(ctx, serverArtist); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddArtist(ctx, serverArtist); err != nil {
				return 0, 0, err
			}
			created += 1
		}

		for _, size := range []imagemagick.ImageSize{imagemagick.Size320, imagemagick.Size640, imagemagick.Size1024, imagemagick.Size1600, imagemagick.Size2400} {
			filename := filepath.Join(s.imgDir, "artists", aID, fmt.Sprintf("%d.jpg", size))

			if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
				return 0, 0, err
			}

			dst, err := os.Create(filename)
			if err != nil {
				return 0, 0, err
			}

			s.log.DebugContext(ctx, fmt.Sprintf("downloading artist image '%s' to '%s'", aID, filename))

			if err = s.sc.DownloadArtistImage(ctx, aID, size, dst); err != nil {
				serr, ok := err.(serverclient.ErrStatus)
				if !ok || serr.Status >= 500 {
					dst.Close()
					return 0, 0, err
				}
			}

			if err = dst.Close(); err != nil {
				return 0, 0, err
			}
		}

	}

	return created, updated, nil
}

func (s Syncer) syncAlbums(ctx context.Context, tx database.DatabaseFace, albumIDs []string) (created, updated int, err error) {
	for _, aID := range albumIDs {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing album '%s'", aID))
		exists, err := tx.AlbumExists(ctx, aID)
		if err != nil {
			return 0, 0, err
		}

		serverAlbum, err := s.sc.GetAlbum(ctx, aID)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			if err = tx.UpdateAlbum(ctx, serverAlbum); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddAlbum(ctx, serverAlbum); err != nil {
				return 0, 0, err
			}
			created += 1
		}

		for _, size := range []imagemagick.ImageSize{imagemagick.Size320, imagemagick.Size640, imagemagick.Size1024, imagemagick.Size1600, imagemagick.Size2400} {
			filename := filepath.Join(s.imgDir, "albums", aID, fmt.Sprintf("%d.jpg", size))

			if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
				return 0, 0, err
			}

			dst, err := os.Create(filename)
			if err != nil {
				return 0, 0, err
			}

			s.log.DebugContext(ctx, fmt.Sprintf("downloading album cover '%s' to '%s'", aID, filename))

			if err = s.sc.DownloadAlbumCover(ctx, aID, size, dst); err != nil {
				serr, ok := err.(serverclient.ErrStatus)
				if !ok || serr.Status >= 500 {
					dst.Close()
					return 0, 0, err
				}
			}

			if err = dst.Close(); err != nil {
				return 0, 0, err
			}
		}

	}

	return created, updated, nil
}

func (s Syncer) syncTracks(ctx context.Context, tx database.DatabaseFace, trackIDs []string) (created, updated int, err error) {
	for _, tID := range trackIDs {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing track '%s'", tID))
		exists, err := tx.TrackExists(ctx, tID)
		if err != nil {
			return 0, 0, err
		}

		serverTrack, err := s.sc.GetTrack(ctx, tID)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			if err = tx.UpdateTrack(ctx, serverTrack); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddTrack(ctx, serverTrack); err != nil {
				return 0, 0, err
			}
			created += 1
		}
	}

	return created, updated, nil
}

func (s Syncer) syncAlbumArtists(ctx context.Context, tx database.DatabaseFace, albumArtists []uuid.UUID, deletedAlbumArtists []uuid.UUID) (created, updated int, err error) {
	for _, aa := range deletedAlbumArtists {
		s.log.DebugContext(ctx, fmt.Sprintf("deleting album artist '%s'", aa.String()))
		if err := tx.DeleteAlbumArtist(ctx, aa); err != nil {
			return 0, 0, err
		}
	}

	for _, aa := range albumArtists {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing album artist '%s'", aa.String()))
		exists, err := tx.AlbumArtistExistsByID(ctx, aa)
		if err != nil {
			return 0, 0, err
		}

		albumArtist, err := s.sc.GetAlbumArtistByID(ctx, aa)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			if err = tx.UpdateAlbumArtist(ctx, albumArtist); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddAlbumArtist(ctx, albumArtist); err != nil {
				return 0, 0, err
			}
			created += 1
		}
	}

	return created, updated, nil
}

func (s Syncer) syncAlbumTracks(ctx context.Context, tx database.DatabaseFace, albumTracks []uuid.UUID) (created, updated int, err error) {
	for _, at := range albumTracks {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing album track '%s'", at.String()))
		exists, err := tx.AlbumTrackExistsByID(ctx, at)
		if err != nil {
			return 0, 0, err
		}

		albumTrack, err := s.sc.GetAlbumTrackByID(ctx, at)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			if err = tx.UpdateAlbumTrack(ctx, albumTrack); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddAlbumTrack(ctx, albumTrack); err != nil {
				return 0, 0, err
			}
			created += 1
		}
	}

	return created, updated, nil
}

func (s Syncer) syncPlaylistTracks(ctx context.Context, tx database.DatabaseFace, added []uuid.UUID, deleted []uuid.UUID) (created, updated int, err error) {
	// user, err := auth.UserFromContext(ctx)
	// if err != nil {
	// 	return 0, 0, err
	// }

	for _, d := range deleted {
		s.log.DebugContext(ctx, fmt.Sprintf("deleting playlist_track '%s'", d.String()))
		if err := tx.DeletePlaylistTrack(ctx, d); err != nil {
			return 0, 0, err
		}
	}

	for _, a := range added {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing playlist '%s'", a.String()))
		exists, err := tx.PlaylistTrackExists(ctx, a)
		if err != nil {
			return 0, 0, err
		}

		pt, err := s.sc.GetPlaylistTrack(ctx, a)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			// FIXME
			// if err = tx.UpdatePlaylistTrack(ctx, playlist); err != nil {
			// 	return 0, 0, err
			// }
			// updated += 1
		} else {
			if _, err = tx.AddPlaylistTrack(ctx, pt); err != nil {
				return 0, 0, err
			}
			created += 1
		}

	}

	return created, updated, nil
}

func (s Syncer) syncPlaylists(ctx context.Context, tx database.DatabaseFace, added []uuid.UUID, deleted []uuid.UUID) (created, updated int, err error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return 0, 0, err
	}

	for _, d := range deleted {
		s.log.DebugContext(ctx, fmt.Sprintf("deleting playlist '%s'", d.String()))
		if err := tx.DeletePlaylist(ctx, d, user.ID); err != nil {
			return 0, 0, err
		}
	}

	for _, a := range added {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing playlist '%s'", a.String()))
		exists, err := tx.PlaylistExistsByID(ctx, a)
		if err != nil {
			return 0, 0, err
		}

		playlist, err := s.sc.GetPlaylist(ctx, a)
		if err != nil {
			return 0, 0, err
		}

		if exists {
			if err = tx.UpdatePlaylist(ctx, playlist); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddPlaylist(ctx, playlist); err != nil {
				return 0, 0, err
			}
			created += 1
		}

		for _, size := range []imagemagick.ImageSize{imagemagick.Size320, imagemagick.Size640, imagemagick.Size1024, imagemagick.Size1600, imagemagick.Size2400} {
			filename := filepath.Join(s.imgDir, "playlists", a.String(), fmt.Sprintf("%d.jpg", size))

			if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
				return 0, 0, err
			}

			dst, err := os.Create(filename)
			if err != nil {
				return 0, 0, err
			}

			s.log.DebugContext(ctx, fmt.Sprintf("downloading playlist image '%s' to '%s'", a.String(), filename))

			if err = s.sc.DownloadPlaylistImage(ctx, a, size, dst); err != nil {
				serr, ok := err.(serverclient.ErrStatus)
				if !ok || serr.Status >= 500 {
					dst.Close()
					return 0, 0, err
				}
			}

			if err = dst.Close(); err != nil {
				return 0, 0, err
			}
		}
	}

	return created, updated, nil
}

func (s Syncer) syncMediaFiles(ctx context.Context, tx database.DatabaseFace, mediaFiles []uuid.UUID) error {
	for _, mfID := range mediaFiles {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing media file '%s'", mfID.String()))

		_, err := tx.AddSyncItem(ctx, types.SyncItem{
			// FIXME: This should be added
			SyncID: uuid.Nil,
			ItemID: mfID,
			Type:   types.SiTypeMediaFile,
			State:  types.SiStateNotStarted,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Syncer) syncPlayHistory(ctx context.Context, tx database.DatabaseFace, lastSync time.Time) error {
	clientHistory, err := tx.ListUpdatedPlayHistory(ctx, lastSync)
	if err != nil {
		return err
	}

	serverSyncState, err := s.sc.SyncPlayHistory(ctx, lastSync, clientHistory)
	if err != nil {
		return err
	}

	for _, serverItem := range serverSyncState.RemoteItems {
		if _, err := tx.AddPlayHistoryStruct(ctx, serverItem); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) SetDatabase(db database.DatabaseFace) {
	s.db = db
}

func New(sc *serverclient.ServerClient, db database.DatabaseFace, musicDir, imgDir string, slogHandler slog.Handler) Syncer {
	return Syncer{
		sc:           sc,
		db:           db,
		musicDir:     musicDir,
		imgDir:       imgDir,
		serverStatus: server.Status{Status: server.HealthzUnavailable},
		log:          slog.New(slogHandler).With("service", "syncer"),
	}
}
