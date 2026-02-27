package server

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type ReqArtistsGet struct {
	ArtistID uuid.UUID `path:"artistID" description:"ID of the wanted artist"`
}

func (r ReqArtistsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqArtistsUpdate struct {
	ArtistID uuid.UUID `path:"artistID" description:"ID of the wanted artist"`
	types.ArtistBase
}

func (r ReqArtistsUpdate) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqList struct {
	Count bool `query:"count" description:"Only return the count, not the payload."`
}

func (r ReqList) Validate() (validationMessages string, err error) {
	return "", nil
}
