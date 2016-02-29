package xerrors

import (
	"fmt"
	"net/http"
)

func PrintStatusInternalServerError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "%v", err.Error())
}
