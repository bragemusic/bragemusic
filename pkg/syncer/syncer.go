package syncer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/bragemusic/core/pkg/utils"
	"github.com/dhowden/tag"
)

type Syncer struct {
	sc       *serverclient.ServerClient
	db       database.DatabaseFace
	log      *slog.Logger
	musicDir string
}

func (s Syncer) Sync(ctx context.Context) error {
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
				// FIXME: This entire thing can be replaced with the actual FilePath from the serverTrack when it is a relative path in the server db
				album, err := tx.GetAlbumFromID(ctx, *serverTrack.AlbumID)
				if err != nil {
					return 0, 0, err
				}

				artist, err := tx.GetArtistFromID(ctx, album.ArtistID)
				if err != nil {
					return 0, 0, err
				}

				albumFolder := utils.GenerateAlbumFolderPath(artist.Name, album.Name, s.musicDir)

				if err = os.MkdirAll(albumFolder, os.ModePerm); err != nil {
					return 0, 0, err
				}

				trackPath, err := utils.GenerateTrackPath(*serverTrack.DiscNumber, *serverTrack.TrackNumber, serverTrack.Title, tag.FileType(*serverTrack.MimeType), albumFolder)
				if err != nil {
					return 0, 0, err
				}

				// trackPath := serverTrack.FilePath
				//////////////////////////////////////////////////////////////////////

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

func New(sc *serverclient.ServerClient, db database.DatabaseFace, musicDir string, slogHandler slog.Handler) Syncer {
	return Syncer{
		sc:       sc,
		db:       db,
		musicDir: musicDir,
		log:      slog.New(slogHandler).With("service", "syncer"),
	}
}
