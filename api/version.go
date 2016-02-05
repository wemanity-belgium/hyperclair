package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/clair"
)

const Version = "1"

type version struct {
	APIVersion string
	Clair      interface{}
}

func (version version) asJSON() string {
	b, err := json.Marshal(version)
	if err != nil {
		fmt.Println(err)
		return string("Cannot marshal health")
	}
	return string(b)
}

func VersionsHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	clairVersion, err := clair.Versions()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error: %v", err)
		return
	}

	versionBody := version{
		APIVersion: Version,
		Clair:      clairVersion,
	}

	fmt.Fprint(rw, versionBody.asJSON())
}
