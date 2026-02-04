package bragerr

import (
	"fmt"
	"net/http"
)

var ErrDepMissing = BragErr{
	Code:   "DEPENDENCYMISSING",
	Title:  "Dependency Missing",
	Status: http.StatusInternalServerError,
}

func (b BragErrFactory) DependencyMissing(err error, dependencyName string) *BragErr {
	e := ErrDepMissing
	e.Service = b.service
	e.Err = err
	e.Message = fmt.Sprintf("Dependency '%s' is missing", dependencyName)
	return &e
}
