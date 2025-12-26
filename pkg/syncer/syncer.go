package syncer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/server"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
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

	if h.Status != server.HealthzRunning {
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

	dbSyncState.ArtistsCreated, dbSyncState.ArtistsUpdated, err = s.syncArtists(ctx, tx, syncState.Artists)
	if err != nil {
		return err
	}

	dbSyncState.AlbumsCreated, dbSyncState.AlbumsUpdated, err = s.syncAlbums(ctx, tx, syncState.Albums)
	if err != nil {
		return err
	}

	dbSyncState.TracksCreated, dbSyncState.TracksUpdated, err = s.syncTracks(ctx, tx, syncState.Tracks)
	if err != nil {
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

		filename := filepath.Join(s.imgDir, "artists", fmt.Sprintf("%s.jpg", aID))

		if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
			return 0, 0, err
		}

		dst, err := os.Create(filename)
		if err != nil {
			return 0, 0, err
		}

		s.log.DebugContext(ctx, fmt.Sprintf("downloading artist image '%s' to '%s'", aID, filename))

		if err = s.sc.DownloadArtistImage(ctx, aID, dst); err != nil {
			dst.Close()
			return 0, 0, err
		}

		if err = dst.Close(); err != nil {
			return 0, 0, err
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

		filename := filepath.Join(s.imgDir, "albums", fmt.Sprintf("%s.jpg", aID))

		if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
			return 0, 0, err
		}

		dst, err := os.Create(filename)
		if err != nil {
			return 0, 0, err
		}

		s.log.DebugContext(ctx, fmt.Sprintf("downloading album cover '%s' to '%s'", aID, filename))

		if err = s.sc.DownloadAlbumCover(ctx, aID, dst); err != nil {
			dst.Close()
			return 0, 0, err
		}

		if err = dst.Close(); err != nil {
			return 0, 0, err
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
			// FIXME: We need a file_updated_at field. So we dont download files too many times
			if err = tx.UpdateTrack(ctx, serverTrack); err != nil {
				return 0, 0, err
			}
			updated += 1
		} else {
			if _, err = tx.AddTrack(ctx, serverTrack); err != nil {
				return 0, 0, err
			}
			created += 1

			if serverTrack.FilePath != "" {
				trackPath := filepath.Join(s.musicDir, serverTrack.FilePath)

				if err = os.MkdirAll(filepath.Dir(trackPath), os.ModePerm); err != nil {
					return 0, 0, err
				}

				dst, err := os.Create(trackPath)
				if err != nil {
					return 0, 0, err
				}

				s.log.InfoContext(ctx, fmt.Sprintf("downloading track '%s' to '%s'", tID, trackPath))

				if err = s.sc.DownloadTrackFile(ctx, tID, dst); err != nil {
					dst.Close()
					return 0, 0, err
				}

				if err = dst.Close(); err != nil {
					return 0, 0, err
				}

			}
		}
	}

	return created, updated, nil
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
