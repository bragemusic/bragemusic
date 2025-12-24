package server

import (
	"fmt"
)

type ServerError interface {
	UserError() string
	Error() string
}

type ErrIDNotFound struct {
	idKey string
	err   error
}

func (e ErrIDNotFound) UserError() string {
	return fmt.Sprintf("%s does not exist", e.idKey)
}

func (e ErrIDNotFound) Error() string {
	return e.err.Error()
}

type ErrUnauthorized struct{}

func (e ErrUnauthorized) UserError() string {
	return "email or password is incorrect"
}

func (e ErrUnauthorized) Error() string {
	return "email or password is incorrect"
}
