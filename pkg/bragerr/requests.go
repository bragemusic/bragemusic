package bragerr

import (
	"errors"
	"net/http"
)

var ErrReqValidation = BragErr{
	Code:   "REQVALIDATION",
	Title:  "Request validation failed",
	Status: http.StatusBadRequest,
}

func (b BragErrFactory) ReqValidation(validationErrs string) *BragErr {
	e := ErrReqValidation
	e.Service = b.service
	e.Err = errors.New("validation of request failed")
	e.Message = validationErrs

	return &e
}
