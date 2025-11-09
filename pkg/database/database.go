package database

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/bragemusic/core/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type DatabaseFace interface {
	Begin(ctx context.Context) (DatabaseFace, error)
	driver.Tx

	AddAlbum(ctx context.Context, a types.Album) (string, error)
	GetAlbumFromArtistAndName(ctx context.Context, artistName, albumName string) (album types.Album, err error)
	GetAlbumFromMbID(ctx context.Context, mbID string) (album types.Album, err error)
	GetAlbumFromID(ctx context.Context, id string) (album types.Album, err error)
	GetEnhancedAlbumFromID(ctx context.Context, id string) (album types.AlbumEnhanced, err error)
	GetAlbumsByMbIDs(ctx context.Context, albumMbIds []string) ([]types.Album, error)
	ListAlbumsByArtist(ctx context.Context, artistID string) (albums []types.Album, err error)

	AddArtist(ctx context.Context, a types.Artist) (string, error)
	UpdateArtist(ctx context.Context, a types.Artist) error
	GetArtistFromID(ctx context.Context, id string) (artist types.Artist, err error)
	GetArtistFromMbID(ctx context.Context, mbID string) (artist types.Artist, err error)
	GetArtistFromName(ctx context.Context, name string) (artist types.Artist, err error)
	ListArtists(ctx context.Context) (artists []types.Artist, err error)

	AddTrack(ctx context.Context, t types.Track) (string, error)
	AddTracks(ctx context.Context, tracks []types.Track) ([]string, error)
	UpdateTrack(ctx context.Context, t types.Track) error
	UpdateTrackFromMbID(ctx context.Context, t types.Track) error
	GetTrackFromMbID(ctx context.Context, mbID string) (track types.Track, err error)
	GetTrackFromID(ctx context.Context, ID string) (track types.Track, err error)
	GetTracksFromAlbumID(ctx context.Context, albumID string) (tracks []types.Track, err error)
	GetEnhancedTracksFromAlbumID(ctx context.Context, albumID string) (tracks []types.TrackEnhanced, err error)
	GetTrackFromName(ctx context.Context, albumID string, trackName string) (track types.Track, err error)

	// // User Handling
	// AddUser(ctx context.Context, user *DbUser) (bool, error)
	// UserExists(ctx context.Context, username string) (bool, error)
	// GetUser(ctx context.Context, username string) (user *DbUser, err error)

	// // Login Session
	// NewLoginSession(ctx context.Context, username string) (dbLoginSession *DbLoginSession, err error)
	// GetLoginSessionByID(ctx context.Context, loginSessionId string) (dbLoginSession *DbLoginSession, err error)
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

	// db.Exec("PRAGMA journal_mode = WAL;")
	// db.Exec("PRAGMA synchronous = NORMAL;") // optional, faster writes
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)

	return Database{
		ext: db,
	}, nil
}
