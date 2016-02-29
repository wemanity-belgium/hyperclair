package api

import (
	"log"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
)

//PushHandler push image to Clair
func PushHandler(rw http.ResponseWriter, request *http.Request) error {
	image, err := docker.Pull(parseImageURL(request))
	if err != nil {
		return err
	}

	log.Println("Pushing Image")
	if err := docker.Push(image); err != nil {
		return err
	}
	rw.WriteHeader(http.StatusCreated)
	return nil
}
