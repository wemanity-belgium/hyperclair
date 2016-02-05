package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
)

func ReportHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	image, err := docker.Parse(parseImageURL(request))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	if err := image.Pull(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	analyses := image.Analyse()
	analysesHTML, err := analyses.ReportAsHTML()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	fmt.Fprint(rw, analysesHTML)
}
