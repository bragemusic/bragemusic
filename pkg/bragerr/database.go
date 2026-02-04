package bragerr

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrDbNotFound = BragErr{
		Code:   "DB_NOT_FOUND",
		Title:  "Not found",
		Status: http.StatusNotFound,
	}

	ErrDbConflict = BragErr{
		Code:   "DB_CONFLICT",
		Title:  "Conflict",
		Status: http.StatusConflict,
	}

	ErrDbInvalidReference = BragErr{
		Code:   "DB_INVALID_REFERENCE",
		Title:  "Invalid reference",
		Status: http.StatusBadRequest,
	}

	ErrDbDatabase = BragErr{
		Code:   "DB_DATABASE_ERROR",
		Title:  "Database error",
		Status: http.StatusInternalServerError,
	}
)

func (b BragErrFactory) DatabaseError(err error, entity types.EntityType, id *uuid.UUID) *BragErr {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		e := ErrDbNotFound

		if id != nil {
			e.Message = fmt.Sprintf(
				"The %s with id '%s' does not exist.",
				entity,
				id.String(),
			)
		} else {
			e.Message = fmt.Sprintf(
				"No %ss were found.",
				entity,
			)
		}

		e.Service = b.service
		e.Err = err
		return &e
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.ExtendedCode {

		case sqlite3.ErrConstraintUnique,
			sqlite3.ErrConstraintPrimaryKey:

			e := ErrDbConflict

			if id != nil {
				e.Message = fmt.Sprintf(
					"A %s with id '%s' already exists.",
					entity,
					id.String(),
				)
			} else {
				e.Message = fmt.Sprintf(
					"A %s with the same identifier already exists.",
					entity,
				)
			}

			e.Service = b.service
			e.Err = err
			return &e

		case sqlite3.ErrConstraintForeignKey:
			e := ErrDbInvalidReference
			e.Message = fmt.Sprintf(
				"The %s references another entity that does not exist.",
				entity,
			)
			e.Service = b.service
			e.Err = err
			return &e
		}
	}

	// 3️⃣ Fallback
	e := ErrDbDatabase
	e.Message = "An unexpected database error occurred."
	e.Service = b.service
	e.Err = err
	return &e
}
