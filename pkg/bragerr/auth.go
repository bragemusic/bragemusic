package bragerr

import (
	"fmt"
	"net/http"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var ErrUnauthenticated = BragErr{
	Code:    "UNAUTHENTICATED",
	Title:   "Authentication required",
	Status:  http.StatusUnauthorized,
	Message: "Valid authentication credentials were not provided.",
}

var ErrItemAccessDenied = BragErr{
	Code:   "FORBIDDEN",
	Title:  "Access Denied",
	Status: http.StatusForbidden,
}

var ErrNoUserInContext = BragErr{
	Code:    "NOUSERINCONTEXT",
	Title:   "No user in context",
	Message: "There is no user logged in",
	Status:  http.StatusForbidden,
}

var ErrInvalidUserCreds = BragErr{
	Code:    "INVALIDUSERCREDS",
	Title:   "Invalid user credentials",
	Status:  http.StatusUnauthorized,
	Message: "Username and/or password is incorrect.",
}

func (b BragErrFactory) Unauthenticated(err error) *BragErr {
	e := ErrUnauthenticated
	e.Service = b.service
	e.Err = err
	return &e
}

func (b BragErrFactory) ItemAccessDenied(err error, entityType types.EntityType, id uuid.UUID) *BragErr {
	e := ErrItemAccessDenied
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("User does not have access to the %s %s", entityType, id)
	return &e
}

func (b BragErrFactory) NoUserInContext(err error) *BragErr {
	e := ErrNoUserInContext
	e.Service = b.service
	e.Err = err
	return &e
}

func (b BragErrFactory) ErrInvalidUserCreds(err error) *BragErr {
	e := ErrInvalidUserCreds
	e.Service = b.service
	e.Err = err
	return &e
}
