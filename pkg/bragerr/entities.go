package bragerr

import (
	"fmt"
	"net/http"

	"github.com/bragemusic/core/pkg/types"
)

var ErrEntityAlreadyExists = BragErr{
	Code:   "ENTITYALREADYEXISTS",
	Title:  "Entity already exists",
	Status: http.StatusConflict,
}

var ErrEntityDoNotExist = BragErr{
	Code:   "ENTITYDONOTEXIST",
	Title:  "Entity do not exist",
	Status: http.StatusNotFound,
}

func (b BragErrFactory) EntityAlreadyExists(err error, entityType types.EntityType) *BragErr {
	e := ErrEntityAlreadyExists
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("cannot create another %s", entityType)

	return &e
}

func (b BragErrFactory) EntityDoNotExist(err error, entityType types.EntityType) *BragErr {
	e := ErrEntityDoNotExist
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("requested %s was not found", entityType)

	return &e
}
