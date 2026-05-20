package bragerr

import (
	"fmt"
	"net/http"
)

var ErrNoMediaFile = BragErr{
	Code:   "MEDIA_NO_MEDIA_FILE",
	Title:  "Not found",
	Status: http.StatusNotFound,
}

func (b BragErrFactory) NoMediaFile(err error, trackName string) *BragErr {
	e := ErrNoMediaFile
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("Track '%s' has no linked media file.", trackName)
	return &e
}
