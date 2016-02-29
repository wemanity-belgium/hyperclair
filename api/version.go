package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/clair"
)

const v = "1"

func VersionsHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")
	clairVersion, err := clair.Versions()

	if err != nil {
		return err
	}

	version := struct {
		APIVersion string
		Clair      interface{}
	}{v, clairVersion}

	b, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("cannot marshal version:%v", err)
	}
	fmt.Fprint(rw, string(b))
	return nil
}
