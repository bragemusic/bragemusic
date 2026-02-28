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

type ReqAlbumsGet struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
}

func (r ReqAlbumsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsTrackGet struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	TrackID uuid.UUID `path:"trackID" description:"ID of the wanted track"`
}

func (r ReqAlbumsTrackGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsAlbumArtistGet struct {
	AlbumID  uuid.UUID        `path:"albumID" description:"ID of the wanted album"`
	ArtistID uuid.UUID        `path:"artistID" description:"ID of the wanted artist"`
	Role     types.ArtistRole `path:"role" description:"Role of the wanted artist"`
}

func (r ReqAlbumsAlbumArtistGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsAlbumTrackGet struct {
	AlbumID     uuid.UUID `path:"albumID" description:"ID of the wanted album track"`
	DiscNumber  int       `path:"discNumber" description:"Disc number of the wanted album track"`
	TrackNumber int       `path:"trackNumber" description:"Track number of the wanted album track"`
}

func (r ReqAlbumsAlbumTrackGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumsUpdate struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	types.AlbumUpdate
}

func (r ReqAlbumsUpdate) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumArtistsGet struct {
	AlbumArtistID uuid.UUID `path:"albumArtistID" description:"ID of the wanted album artist"`
}

func (r ReqAlbumArtistsGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqAlbumTracksGet struct {
	AlbumTrackID uuid.UUID `path:"albumTrackID" description:"ID of the wanted album track"`
}

func (r ReqAlbumTracksGet) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqList struct {
	Count bool `query:"count" description:"Only return the count, not the payload."`
}

func (r ReqList) Validate() (validationMessages string, err error) {
	return "", nil
}

type ReqListTracksOfAlbum struct {
	AlbumID uuid.UUID `path:"albumID" description:"ID of the wanted album"`
	Count   bool      `query:"count" description:"Only return the count, not the payload."`
}

func (r ReqListTracksOfAlbum) Validate() (validationMessages string, err error) {
	return "", nil
}
