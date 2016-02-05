package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
)

func AnalyseHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
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
	analysesJSON, err := analyses.ReportAsJSON()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error: %v", err)
	}

	fmt.Fprint(rw, string(analysesJSON))
}
