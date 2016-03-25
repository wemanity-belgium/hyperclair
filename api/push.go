package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/docker"
)

//PushHandler push image to Clair
func PushHandler(rw http.ResponseWriter, request *http.Request) error {
	local := request.URL.Query()["local"]

	docker.IsLocal = len(local) > 0

	image, err := docker.Parse(parseImageURL(request))
	if !docker.IsLocal {
		image, err = docker.Pull(parseImageURL(request))
		if err != nil {
			return err
		}
	}

	logrus.Info("Pushing Image")
	if err := docker.Push(image); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusCreated)
	return nil
}
