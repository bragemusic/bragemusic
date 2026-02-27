package server

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type Response struct {
	Payload any `json:"payload,omitempty"`
	Status  int `json:"-"`
}

type Request[T any] struct{}

func (r Request[T]) Validate() (validationMessages string, err error) {
	return "", nil
}

type Request1 struct {
	AlbumID    uuid.UUID `path:"albumID"`
	SearchTerm string    `query:"q" required:"true" description:"Filter the album names"`
	types.RatingReq
}

func (r Request1) Validate() (validationMessages string, err error) {
	return "", nil
}
