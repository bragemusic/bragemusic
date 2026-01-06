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

var ErrNoRowDeleted = errors.New("no row was deleted")

type DatabaseFace interface {
	Begin(ctx context.Context) (DatabaseFace, error)
	driver.Tx

	AddAlbum(ctx context.Context, a types.Album) (uuid.UUID, error)
	AlbumExists(ctx context.Context, ID string) (bool, error)
	UpdateAlbum(ctx context.Context, a types.Album) error
	GetAlbumFromArtistAndName(ctx context.Context, artistName, albumName string) (album types.Album, err error)
	GetAlbumFromMbID(ctx context.Context, mbID string) (album types.Album, err error)
	GetAlbumFromID(ctx context.Context, id string) (album types.Album, err error)
	GetAlbumDetailed(ctx context.Context, albumID uuid.UUID) (album types.AlbumDetailed, err error)
	GetAlbumsByMbIDs(ctx context.Context, albumMbIds []string) ([]types.Album, error)
	ListAlbumsByArtist(ctx context.Context, artistID string, sortBy SortBy, sortOrder SortOrder) (albums []types.AlbumDetailed, err error)
	ListUpdatedAlbums(ctx context.Context, since time.Time) (albumIDs []string, err error)

	AddArtist(ctx context.Context, a types.Artist) (uuid.UUID, error)
	ArtistExists(ctx context.Context, ID string) (bool, error)
	UpdateArtist(ctx context.Context, a types.Artist) error
	GetArtistFromID(ctx context.Context, id string) (artist types.Artist, err error)
	GetArtistFromMbID(ctx context.Context, mbID string) (artist types.Artist, err error)
	GetArtistFromName(ctx context.Context, name string) (artist types.Artist, err error)
	ListArtists(ctx context.Context, sortBy SortBy, sortOrder SortOrder) (artists []types.Artist, err error)
	ListUpdatedArtists(ctx context.Context, since time.Time) (artistIDs []string, err error)
	ListArtistsWithoutMeta(ctx context.Context) (artists []types.Artist, err error)

	AddTrack(ctx context.Context, t types.Track) (uuid.UUID, error)
	TrackExists(ctx context.Context, ID string) (bool, error)
	TrackExistsByNameAndAlbumID(ctx context.Context, title, albumID string) (bool, error)
	UpdateTrack(ctx context.Context, t types.Track) error
	UpdateTrackFromMbID(ctx context.Context, t types.Track) error
	GetTrackFromMbID(ctx context.Context, mbID string) (track types.Track, err error)
	GetTrackFromID(ctx context.Context, ID string) (track types.Track, err error)
	GetTracksFromAlbumID(ctx context.Context, albumID string) (tracks []types.Track, err error)
	GetTracksDetailedFromArtistID(ctx context.Context, artistID string, sortBy SortBy, sortOrder SortOrder, limit *int, includeMissingFiles bool) (tracks []types.TrackDetailed, err error)
	ListAlbumTracksDetailed(ctx context.Context, albumID uuid.UUID) (tracks []types.TrackDetailed, err error)
	GetTrackFromName(ctx context.Context, albumID uuid.UUID, trackName string) (track types.Track, err error)
	ListUpdatedTracks(ctx context.Context, since time.Time) (trackIDs []string, err error)

	AddMediaFile(ctx context.Context, mf types.MediaFile) (uuid.UUID, error)
	GetMediaFileFromChecksum(ctx context.Context, cs string) (mf types.MediaFile, err error)
	ListUpdatedMediaFiles(ctx context.Context, since time.Time) (mediaFileIDs []uuid.UUID, err error)

	AddAlbumTrack(ctx context.Context, at types.AlbumTrack) error
	AlbumTrackExists(ctx context.Context, albumID uuid.UUID, trackID uuid.UUID) (bool, error)
	ListUpdatedAlbumTracks(ctx context.Context, since time.Time) (albumTracks []types.AlbumTrackKey, err error)

	AddAlbumArtist(ctx context.Context, aa types.AlbumArtist) error
	AlbumArtistExists(ctx context.Context, albumID uuid.UUID, artistID uuid.UUID, role types.ArtistRole) (bool, error)
	ListUpdatedAlbumArtists(ctx context.Context, since time.Time) (albumArtists []types.AlbumArtistKey, err error)
	GetAlbumArtist(ctx context.Context, albumID, artistID uuid.UUID, role types.ArtistRole) (albumArtist types.AlbumArtist, err error)
	UpdateAlbumArtist(ctx context.Context, aa types.AlbumArtist) error

	AddSync(ctx context.Context, s types.DBSyncState) (string, error)
	GetLastSync(ctx context.Context) (sync types.DBSyncState, err error)

	AddPlayHistory(ctx context.Context, trackID, userID uuid.UUID) (string, error)
	AddPlayHistoryStruct(ctx context.Context, ph types.PlayHistory) (string, error)
	ListUpdatedPlayHistory(ctx context.Context, since time.Time) (updatedItems []types.PlayHistory, err error)

	AuthFace
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
		return conn.RegisterFunc("normalize", normalizeForCompare, true)
	}

	// Allows for cascade deleting of foreign rows (user related at the moment)
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return Database{}, err
	}
	// db.Exec("PRAGMA journal_mode = WAL;")
	// db.Exec("PRAGMA synchronous = NORMAL;") // optional, faster writes
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)

	return Database{
		ext: db,
	}, nil
}
