package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/docker"
)

func ReportHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "text/html")
	image, err := docker.Pull(parseImageURL(request))

	if err != nil {
		return err
	}

	analyses := docker.Analyse(image)
	analysesHTML, err := clair.ReportAsHTML(analyses)
	if err != nil {
		return err
	}

	fmt.Fprint(rw, analysesHTML)
	return nil
}
