package database

import (
	"context"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type (
	SortBy    string
	SortOrder string
)

const (
	SortByName      SortBy = "name"
	SortByDate      SortBy = "date"
	SortByPlayCount SortBy = "play_count"

	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)

var (
	ErrNoRowDeleted = errors.New("no row was deleted")
	ErrNoUser       = errors.New("no user id provided")
)

type DatabaseFace interface {
	Begin(ctx context.Context) (DatabaseFace, error)
	driver.Tx

	AddAlbum(ctx context.Context, a types.Album, userID uuid.UUID) (uuid.UUID, error)
	AlbumExists(ctx context.Context, ID string) (bool, error)
	UpdateAlbum(ctx context.Context, a types.Album, userID uuid.UUID) error
	GetAlbumFromArtistAndName(ctx context.Context, artistName, albumName string) (album types.Album, err error)
	GetAlbumFromMbID(ctx context.Context, mbID string) (album types.Album, err error)
	GetAlbumFromID(ctx context.Context, id uuid.UUID) (album types.Album, err error)
	GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (album types.AlbumDetailed, err error)
	GetAlbumsByMbIDs(ctx context.Context, albumMbIds []string) ([]types.Album, error)
	ListAlbums(ctx context.Context) (albums []types.Album, err error)
	ListAlbumsByArtist(ctx context.Context, artistID uuid.UUID, sortBy SortBy, sortOrder SortOrder) (albums []types.AlbumDetailed, err error)
	ListUpdatedAlbums(ctx context.Context, since time.Time) (albumIDs []string, err error)
	CountAlbums(ctx context.Context) (int, error)
	CountAlbumsByArtist(ctx context.Context, artistID uuid.UUID) (int, error)

	AddArtist(ctx context.Context, a types.Artist, userID uuid.UUID) (uuid.UUID, error)
	ArtistExists(ctx context.Context, ID string) (bool, error)
	UpdateArtist(ctx context.Context, a types.Artist, userID uuid.UUID) error
	GetArtistFromID(ctx context.Context, id uuid.UUID) (artist types.Artist, err error)
	GetArtistFromMbID(ctx context.Context, mbID string) (artist types.Artist, err error)
	GetArtistFromName(ctx context.Context, name string) (artist types.Artist, err error)
	ListArtists(ctx context.Context, sortBy SortBy, sortOrder SortOrder) (artists []types.ArtistDetailed, err error)
	ListUpdatedArtists(ctx context.Context, since time.Time) (artistIDs []string, err error)
	ListArtistsWithoutMeta(ctx context.Context) (artists []types.Artist, err error)
	CountArtists(ctx context.Context) (int, error)

	AddTrack(ctx context.Context, t types.Track, userID uuid.UUID) (uuid.UUID, error)
	TrackExists(ctx context.Context, ID string) (bool, error)
	TrackExistsByNameAndAlbumID(ctx context.Context, title, albumID string) (bool, error)
	UpdateTrack(ctx context.Context, t types.Track, userID uuid.UUID) error
	GetTrackFromMbID(ctx context.Context, mbID string) (track types.Track, err error)
	GetTrackFromID(ctx context.Context, ID uuid.UUID) (track types.Track, err error)
	GetTrackDetailed(ctx context.Context, trackID, albumID, userID uuid.UUID) (track types.TrackDetailed, err error)
	GetTracksFromAlbumID(ctx context.Context, albumID uuid.UUID) (tracks []types.Track, err error)
	GetTracksDetailedFromArtistID(ctx context.Context, artistID, userID uuid.UUID, sortBy SortBy, sortOrder SortOrder, limit *int, includeMissingFiles bool) (tracks []types.TrackDetailed, err error)
	ListTracks(ctx context.Context) (tracks []types.Track, err error)
	ListAlbumTracksDetailed(ctx context.Context, albumID, userID uuid.UUID) (tracks []types.TrackDetailed, err error)
	GetTrackFromName(ctx context.Context, albumID uuid.UUID, trackName string) (track types.Track, err error)
	ListUpdatedTracks(ctx context.Context, since time.Time) (trackIDs []string, err error)
	CountTracks(ctx context.Context) (int, error)

	AddMediaFile(ctx context.Context, mf types.MediaFile, userID uuid.UUID) (uuid.UUID, error)
	GetMediaFile(ctx context.Context, id uuid.UUID) (mf types.MediaFile, err error)
	GetMediaFileFromChecksum(ctx context.Context, cs string) (mf types.MediaFile, err error)
	ListUpdatedMediaFiles(ctx context.Context, since time.Time) (mediaFileIDs []uuid.UUID, err error)
	MediaFileExists(ctx context.Context, ID uuid.UUID) (bool, error)
	UpdateMediaFile(ctx context.Context, mf types.MediaFile, userID uuid.UUID) error

	AddAlbumTrack(ctx context.Context, at types.AlbumTrack, userID uuid.UUID) (uuid.UUID, error)
	AlbumTrackExists(ctx context.Context, albumID uuid.UUID, trackID uuid.UUID) (bool, error)
	AlbumTrackExistsByPos(ctx context.Context, albumID uuid.UUID, discNumber, trackNumber int) (bool, error)
	AlbumTrackExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	ListUpdatedAlbumTracks(ctx context.Context, since time.Time) (albumTracks []uuid.UUID, err error)
	GetAlbumTrack(ctx context.Context, albumID uuid.UUID, discNumber, trackNumber int) (albumTrack types.AlbumTrack, err error)
	GetAlbumTrackByID(ctx context.Context, id uuid.UUID) (albumTrack types.AlbumTrack, err error)
	GetAlbumTrackFromAlbumAndTrack(ctx context.Context, albumID, trackID uuid.UUID) (albumTrack types.AlbumTrack, err error)
	ListAlbumTracksByTrackID(ctx context.Context, trackID uuid.UUID) (results []types.AlbumTrack, err error)
	UpdateAlbumTrack(ctx context.Context, at types.AlbumTrack, userID uuid.UUID) error
	CountTracksByArtist(ctx context.Context, artistID uuid.UUID) (int, error)

	AddAlbumArtist(ctx context.Context, aa types.AlbumArtist, userID uuid.UUID) (uuid.UUID, error)
	AlbumArtistExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	AlbumArtistExists(ctx context.Context, id uuid.UUID, artistID uuid.UUID, role types.ArtistRole) (bool, error)
	ListUpdatedAlbumArtists(ctx context.Context, since time.Time) (albumArtists []uuid.UUID, err error)
	GetAlbumArtist(ctx context.Context, albumID, artistID uuid.UUID, role types.ArtistRole) (albumArtist types.AlbumArtist, err error)
	GetAlbumArtistByID(ctx context.Context, id uuid.UUID) (albumArtist types.AlbumArtist, err error)
	UpdateAlbumArtist(ctx context.Context, aa types.AlbumArtist, userID uuid.UUID) error
	ListAlbumArtistsByAlbumID(ctx context.Context, albumID uuid.UUID) (albumArtists []types.AlbumArtist, err error)
	DeleteAlbumArtist(ctx context.Context, id, userID uuid.UUID) error

	AddSync(ctx context.Context, s types.DBSyncState) (string, error)
	GetLastSync(ctx context.Context) (sync types.DBSyncState, err error)

	AddSyncItem(ctx context.Context, s types.SyncItem) (uuid.UUID, error)
	GetUnsyncedItem(ctx context.Context) (si types.SyncItem, err error)
	SetSyncItemState(ctx context.Context, id uuid.UUID, state types.SyncItemState) error

	AddPlayHistory(ctx context.Context, trackID, userID uuid.UUID) (string, error)
	AddPlayHistoryStruct(ctx context.Context, ph types.PlayHistory) (string, error)
	ListUpdatedPlayHistory(ctx context.Context, since time.Time) (updatedItems []types.PlayHistory, err error)

	ListEntityEventsByType(ctx context.Context, eventType types.EntityEventType, entityType types.EntityType, since time.Time) (ids []uuid.UUID, err error)
	ListEntityEventsByEntityType(ctx context.Context, entityType types.EntityType, since time.Time, userID *uuid.UUID) (ids []types.EntityEvent, err error)

	AddSearchItem(ctx context.Context, si types.SearchItem) error
	DeleteAllSearchItems(ctx context.Context) error
	SearchFull(ctx context.Context, searchTerm string, limit int) (results []types.SearchItem, err error)

	AddPlaylist(ctx context.Context, p types.Playlist, userID uuid.UUID) (uuid.UUID, error)
	AddPlaylistTrack(ctx context.Context, p types.PlaylistTrack, userID uuid.UUID) (uuid.UUID, error)
	CountPlaylists(ctx context.Context, userID uuid.UUID) (int, error)
	CountPlaylistTracks(ctx context.Context, playlistID uuid.UUID) (int, error)
	DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error
	DeletePlaylistTrack(ctx context.Context, id, userID uuid.UUID) error
	GetPlaylist(ctx context.Context, ID, userID uuid.UUID) (plist types.Playlist, err error)
	GetPlaylistTrack(ctx context.Context, id uuid.UUID) (plistTrack types.PlaylistTrack, err error)
	GetPlaylistTrackByPlaylistAndAlbumTrack(ctx context.Context, playlistID, albumTrackID uuid.UUID) (plistTrack types.PlaylistTrack, err error)
	ListPlaylists(ctx context.Context, userID uuid.UUID, includePublic bool, sortBy SortBy, sortOrder SortOrder) (playlists []types.Playlist, err error)
	ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) (tracks []types.TrackDetailed, err error)
	ListUpdatedPlaylists(ctx context.Context, since time.Time, userID uuid.UUID) (plists []uuid.UUID, err error)
	ListUpdatedPlaylistTracks(ctx context.Context, since time.Time, userID uuid.UUID) (plists []uuid.UUID, err error)
	PlaylistExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	PlaylistTrackExists(ctx context.Context, id uuid.UUID) (bool, error)
	UpdatePlaylist(ctx context.Context, plist types.Playlist, userID uuid.UUID) error

	AddRating(ctx context.Context, r types.Rating) (uuid.UUID, error)
	GetRating(ctx context.Context, id uuid.UUID) (rating types.Rating, err error)
	GetRatingID(ctx context.Context, trackID, userID uuid.UUID) (id uuid.UUID, found bool, err error)
	GetTrackRatings(ctx context.Context, trackID uuid.UUID) (ratings []types.Rating, err error)
	UpdateRating(ctx context.Context, id uuid.UUID, rating int, userID uuid.UUID) error

	AuthFace
	ImportFace
}

type executor interface {
	sqlx.Execer
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.Queryer
	sqlx.Ext
}

type Database struct {
	ext sqlx.ExtContext
	// db *sqlx.DB
	// executor
}

func (d Database) Begin(ctx context.Context) (DatabaseFace, error) {
	db, ok := d.ext.(*sqlx.DB)
	if !ok {
		return nil, errors.New("cannot start transaction inside another transaction")
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Database{ext: tx}, nil
}

func (d Database) Commit() error {
	if tx, ok := d.ext.(*sqlx.Tx); ok {
		return tx.Commit()
	}
	return nil
}

func (d Database) Rollback() error {
	if tx, ok := d.ext.(*sqlx.Tx); ok {
		return tx.Rollback()
	}
	return nil
}

func New(db *sqlx.DB) (Database, error) {
	sqlite3conn := db.Driver().(*sqlite3.SQLiteDriver)
	sqlite3conn.ConnectHook = func(conn *sqlite3.SQLiteConn) error {
		if err := conn.RegisterFunc("normalize", normalizeForCompare, true); err != nil {
			return err
		}

		// Enable foreign key enforcement on every connection
		if _, err := conn.Exec("PRAGMA foreign_keys = ON;", nil); err != nil {
			return err
		}

		// Enable multiple readers while writing
		if _, err := conn.Exec("PRAGMA journal_mode=WAL;", nil); err != nil {
			return err
		}

		// Enable multiple readers while writing
		if _, err := conn.Exec("PRAGMA synchronous=NORMAL;", nil); err != nil {
			return err
		}

		// Set timeout to 5s
		if _, err := conn.Exec("PRAGMA busy_timeout = 5000;", nil); err != nil {
			return err
		}

		return nil
	}

	// db.Exec("PRAGMA journal_mode = WAL;")
	// db.Exec("PRAGMA synchronous = NORMAL;") // optional, faster writes
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)

	return Database{
		ext: db,
	}, nil
}
