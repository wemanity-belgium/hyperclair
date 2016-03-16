package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
)

//PullHandler return the Light Manifest representation of the docker image
func PullHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")

	image, err := docker.Pull(parseImageURL(request))
	if err != nil {
		return err
	}

	j, err := image.AsJSON()

	if err != nil {
		return err
	}
	fmt.Fprint(rw, j)
	return nil
}
