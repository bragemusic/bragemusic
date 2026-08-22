package bragerr

import (
	"fmt"
	"net/http"

	"github.com/bragemusic/bragemusic/pkg/types"
)

var ErrJobTypeMissing = BragErr{
	Code:   "JOBTYPEMISSING",
	Title:  "Job type missing",
	Status: http.StatusInternalServerError,
}

func (b BragErrFactory) JobTypeMissing(err error, jobType types.JobType) *BragErr {
	e := ErrJobTypeMissing
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("job type '%s' does not exist", jobType)

	return &e
}
