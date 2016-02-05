package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
)

//PushHandler push image to Clair
func PushHandler(rw http.ResponseWriter, request *http.Request) {
	image, err := docker.Parse(parseImageURL(request))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	if err := image.Pull(); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	log.Println("Pushing Image")
	if err := image.Push(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error: %v", err)
	}
	rw.WriteHeader(http.StatusNoContent)
}
