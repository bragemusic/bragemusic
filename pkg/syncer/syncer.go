package syncer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/serverclient"
)

type Syncer struct {
	sc  *serverclient.ServerClient
	db  database.DatabaseFace
	log *slog.Logger
}

func (s Syncer) Sync(ctx context.Context) error {
	s.log.InfoContext(ctx, "starting sync")
	// FIXME: This needs to be in kept in a table
	lastSync := time.Now().AddDate(0, 0, -20)

	syncState, err := s.sc.GetSyncState(ctx, lastSync)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = s.syncArtists(ctx, tx, syncState.Artists); err != nil {
		return err
	}

	if err = s.syncAlbums(ctx, tx, syncState.Albums); err != nil {
		return err
	}

	// FIXME: This needs to be in kept in a table
	lastSync = syncState.Time

	return tx.Commit()
}

func (s Syncer) syncArtists(ctx context.Context, tx database.DatabaseFace, artistIDs []string) error {
	for _, aID := range artistIDs {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing artist '%s'", aID))
		exists, err := tx.ArtistExists(ctx, aID)
		if err != nil {
			return err
		}

		serverArtist, err := s.sc.GetArtist(ctx, aID)
		if err != nil {
			return err
		}

		if exists {
			if err = tx.UpdateArtist(ctx, serverArtist); err != nil {
				return err
			}
		} else {
			if _, err = tx.AddArtist(ctx, serverArtist); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s Syncer) syncAlbums(ctx context.Context, tx database.DatabaseFace, albumIDs []string) error {
	for _, aID := range albumIDs {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing album '%s'", aID))
		exists, err := tx.AlbumExists(ctx, aID)
		if err != nil {
			return err
		}

		serverAlbum, err := s.sc.GetAlbum(ctx, aID)
		if err != nil {
			return err
		}

		if exists {
			if err = tx.UpdateAlbum(ctx, serverAlbum); err != nil {
				return err
			}
		} else {
			if _, err = tx.AddAlbum(ctx, serverAlbum); err != nil {
				return err
			}
		}
	}

	return nil
}

func New(sc *serverclient.ServerClient, db database.DatabaseFace, slogHandler slog.Handler) Syncer {
	return Syncer{
		sc:  sc,
		db:  db,
		log: slog.New(slogHandler).With("service", "syncer"),
	}
}
