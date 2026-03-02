package syncer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/imagemagick"
	"github.com/bragemusic/core/pkg/serverclient"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type Syncer struct {
	sc                      *serverclient.ServerClient
	db                      database.DatabaseFace
	log                     *slog.Logger
	musicDir                string
	imgDir                  string
	serverAvailable         bool
	syncInProgress          bool
	user                    *types.UserDetails
	syncInProgressCallbacks []func(bool)
}

func (s *Syncer) RegisterSyncInProgressCallback(f func(bool)) {
	s.syncInProgressCallbacks = append(s.syncInProgressCallbacks, f)
}

func (s Syncer) SupportsSync() bool {
	return true
}

func (s *Syncer) SetUser(user *types.UserDetails) {
	s.user = user
}

func (s *Syncer) Sync(ctx context.Context) error {
	if s.user == nil {
		return errors.New("no user in context")
	}

	err := s.sync(ctx, s.user.ID)
	if err != nil {
		return err
	}

	err = s.syncItems(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Syncer) sync(ctx context.Context, userID uuid.UUID) error {
	// if !s.serverAvailable {
	// 	return errors.New("server is not available")
	// }

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

	if err = s.syncEntityEvents(ctx, tx, userID, syncState.New); err != nil {
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

func (s Syncer) syncEntityEvents(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, events types.EntityEvents) error {
	for _, e := range events {
		s.log.DebugContext(ctx, fmt.Sprintf("syncing '%s'", e.ItemID), "type", e.EntityType, "action", e.Type)

		if events.LaterDeleteExists(e) {
			s.log.DebugContext(ctx, fmt.Sprintf("skipping '%s'. There is a later delete event.", e.ItemID), "type", e.EntityType, "action", e.Type)
			continue
		}

		var f func(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) error

		switch e.EntityType {
		case types.EntityArtist:
			f = s.syncArtist
		case types.EntityAlbum:
			f = s.syncAlbum
		case types.EntityTrack:
			f = s.syncTrack
		case types.EntityAlbumArtist:
			f = s.syncAlbumArtist
		case types.EntityAlbumTrack:
			f = s.syncAlbumTrack
		case types.EntityPlaylist:
			f = s.syncPlaylist
		case types.EntityPlaylistTrack:
			f = s.syncPlaylistTrack
		case types.EntityMediaFile:
			f = s.syncMediaFile
		case types.EntityRating:
			f = s.syncRating
		default:
			return fmt.Errorf("unsupported entity type '%s'", e.EntityType)
		}

		if err := f(ctx, tx, userID, e); err != nil {
			return err
		}

	}
	return nil
}

func (s Syncer) syncArtist(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var a types.Artist

	if event.Type != types.EntityEventDelete {
		a, err = s.sc.GetArtist(ctx, event.ItemID)
		if err != nil {
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err := tx.AddArtist(ctx, a, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		if err := tx.UpdateArtist(ctx, a, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		return errors.New("'delete' not supported for albums")
	}

	for _, size := range []imagemagick.ImageSize{imagemagick.Size320, imagemagick.Size640, imagemagick.Size1024, imagemagick.Size1600, imagemagick.Size2400} {
		filename := filepath.Join(s.imgDir, "artists", a.ID.String(), fmt.Sprintf("%d.jpg", size))

		if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
			return err
		}

		dst, err := os.Create(filename)
		if err != nil {
			return err
		}

		s.log.DebugContext(ctx, fmt.Sprintf("downloading artist image '%s' to '%s'", a.ID.String(), filename))

		if err = s.sc.DownloadArtistImage(ctx, a.ID.String(), size, dst); err != nil {
			serr, ok := err.(serverclient.ErrStatus)
			if !ok || serr.Status >= 500 {
				dst.Close()
				return err
			}
		}

		if err = dst.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s Syncer) syncAlbum(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var a types.Album

	if event.Type != types.EntityEventDelete {
		a, err = s.sc.GetAlbum(ctx, event.ItemID)
		if err != nil {
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err := tx.AddAlbum(ctx, a, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		if err := tx.UpdateAlbum(ctx, a, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		return errors.New("'delete' not supported for albums")
	}

	for _, size := range []imagemagick.ImageSize{imagemagick.Size320, imagemagick.Size640, imagemagick.Size1024, imagemagick.Size1600, imagemagick.Size2400} {
		filename := filepath.Join(s.imgDir, "albums", a.ID.String(), fmt.Sprintf("%d.jpg", size))

		if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
			return err
		}

		dst, err := os.Create(filename)
		if err != nil {
			return err
		}

		s.log.DebugContext(ctx, fmt.Sprintf("downloading album cover '%s' to '%s'", a.ID.String(), filename))

		if err = s.sc.DownloadAlbumCover(ctx, a.ID, size, dst); err != nil {
			serr, ok := err.(serverclient.ErrStatus)
			if !ok || serr.Status >= 500 {
				dst.Close()
				return err
			}
		}

		if err = dst.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) syncTrack(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var t types.Track

	if event.Type != types.EntityEventDelete {
		t, err = s.sc.GetTrack(ctx, event.ItemID)
		if err != nil {
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err := tx.AddTrack(ctx, t, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		if err := tx.UpdateTrack(ctx, t, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		return errors.New("'delete' not supported for tracks")
	}

	return nil
}

func (s *Syncer) syncAlbumArtist(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var aa types.AlbumArtist

	if event.Type != types.EntityEventDelete {
		aa, err = s.sc.GetAlbumArtistByID(ctx, event.ItemID)
		if err != nil {
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err := tx.AddAlbumArtist(ctx, aa, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		if err := tx.UpdateAlbumArtist(ctx, aa, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		if err := tx.DeleteAlbumArtist(ctx, event.ItemID, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) syncAlbumTrack(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var at types.AlbumTrack

	if event.Type != types.EntityEventDelete {
		at, err = s.sc.GetAlbumTrackByID(ctx, event.ItemID)
		if err != nil {
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err := tx.AddAlbumTrack(ctx, at, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		if err := tx.UpdateAlbumTrack(ctx, at, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		return errors.New("'delete' not supported for album_tracks")
	}

	return nil
}

func (s *Syncer) syncPlaylist(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var p types.Playlist

	if event.Type != types.EntityEventDelete {
		p, err = s.sc.GetPlaylist(ctx, event.ItemID)
		if err != nil {
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err = tx.AddPlaylist(ctx, p, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		if err = tx.UpdatePlaylist(ctx, p, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		if err = tx.DeletePlaylist(ctx, event.ItemID, userID); err != nil {
			return err
		}
	}

	for _, size := range []imagemagick.ImageSize{imagemagick.Size320, imagemagick.Size640, imagemagick.Size1024, imagemagick.Size1600, imagemagick.Size2400} {
		filename := filepath.Join(s.imgDir, "playlists", event.ItemID.String(), fmt.Sprintf("%d.jpg", size))

		if err = os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
			return err
		}

		dst, err := os.Create(filename)
		if err != nil {
			return err
		}

		s.log.DebugContext(ctx, fmt.Sprintf("downloading playlist image '%s' to '%s'", event.ItemID.String(), filename))

		if err = s.sc.DownloadPlaylistImage(ctx, event.ItemID, size, dst); err != nil {
			serr, ok := err.(serverclient.ErrStatus)
			if !ok || serr.Status >= 500 {
				dst.Close()
				return err
			}
		}

		if err = dst.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) syncPlaylistTrack(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	var p types.PlaylistTrack

	if event.Type != types.EntityEventDelete {
		p, err = s.sc.GetPlaylistTrack(ctx, event.ItemID)
		if err != nil {
			serr, ok := err.(serverclient.ErrStatus)
			if ok && serr.Status == http.StatusNotFound {
				s.log.WarnContext(ctx, "playlist track does not exists on server, has probably been cascade deleted. Skipping", "id", event.ItemID)
				return nil
			}
			return err
		}
	}

	switch event.Type {
	case types.EntityEventCreate:
		if _, err = tx.AddPlaylistTrack(ctx, p, userID); err != nil {
			return err
		}
	case types.EntityEventUpdate:
		return errors.New("'update' not supported for ratings")
	case types.EntityEventDelete:
		if err = tx.DeletePlaylistTrack(ctx, event.ItemID, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) syncMediaFile(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) (err error) {
	_, err = tx.AddSyncItem(ctx, types.SyncItem{
		// FIXME: This should be added
		SyncID: uuid.Nil,
		ItemID: event.ItemID,
		Type:   types.SiTypeMediaFile,
		State:  types.SiStateNotStarted,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Syncer) syncRating(ctx context.Context, tx database.DatabaseFace, userID uuid.UUID, event types.EntityEvent) error {
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
		if err := tx.UpdateRating(ctx, r.ID, r.Rating, userID); err != nil {
			return err
		}
	case types.EntityEventDelete:
		return errors.New("'delete' not supported for ratings")
	}
	return nil
}

func (s *Syncer) syncItems(ctx context.Context) error {
	// if !s.serverAvailable {
	// 	return errors.New("server is not available")
	// }

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
			if err = tx.UpdateMediaFile(ctx, mf, s.user.ID); err != nil {
				return err
			}
		} else {
			if _, err = tx.AddMediaFile(ctx, mf, s.user.ID); err != nil {
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
		sc:       sc,
		db:       db,
		musicDir: musicDir,
		imgDir:   imgDir,
		log:      slog.New(slogHandler).With("service", "syncer"),
	}
}
