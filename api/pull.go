package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
)

//PullHandler return the Light Manifest representation of the docker image
func PullHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	image, err := docker.Parse(parseImageURL(request))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	if err := image.Pull(); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	fmt.Fprint(rw, image.String())
}
