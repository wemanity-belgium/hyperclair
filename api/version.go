package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/api/xerrors"
	"github.com/wemanity-belgium/hyperclair/clair"
)

const v = "1"

func VersionsHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	clairVersion, err := clair.Versions()

	if err != nil {
		xerrors.PrintStatusInternalServerError(rw, err)
		return
	}

	version := struct {
		APIVersion string
		Clair      interface{}
	}{v, clairVersion}

	b, err := json.Marshal(version)
	if err != nil {
		xerrors.PrintStatusInternalServerError(rw, fmt.Errorf("cannot marshal version:%v", err))
		return
	}

	fmt.Fprint(rw, string(b))
}
