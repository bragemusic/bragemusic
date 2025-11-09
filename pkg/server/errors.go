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
