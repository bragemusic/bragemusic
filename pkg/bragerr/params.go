package bragerr

import (
	"fmt"
	"net/http"

	"github.com/bragemusic/core/pkg/types"
)

var ErrParamMissing = BragErr{
	Code:   "PARAMMISSING",
	Title:  "Parameter missing",
	Status: http.StatusBadRequest,
}

var ErrParamWrongFormat = BragErr{
	Code:   "PARAMBADFORMAT",
	Title:  "Parameter has the wrong format",
	Status: http.StatusBadRequest,
}

func (b BragErrFactory) ParamWrongFormat(err error, paramName, paramType string) *BragErr {
	e := ErrParamWrongFormat
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("%s must be of type '%s'", paramName, paramType)

	return &e
}

func (b BragErrFactory) ParamMissing(err error, paramName string, entityType *types.EntityType, action *types.Action) *BragErr {
	e := ErrParamMissing
	e.Service = b.service
	e.Err = err
	if entityType != nil && action != nil {
		e.Message = fmt.Sprintf("Cannot %s %s, parameter %s is missing", *action, *entityType, paramName)
	} else {
		e.Message = fmt.Sprintf("Parameter %s is missing", paramName)
	}
	return &e
}
