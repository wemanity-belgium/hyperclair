package xerrors

import (
	"errors"
	"fmt"
	"net/http"
)

var ServiceUnavailable = errors.New("service is unavailable")

func PrintStatusInternalServerError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "%v", err.Error())
}
