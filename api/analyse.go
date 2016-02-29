package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

func AnalyseHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")

	image, err := docker.Pull(parseImageURL(request))

	if err != nil {
		return err
	}

	analyses := docker.Analyse(image)
	analysesJSON, err := xstrings.ToIndentJSON(analyses)
	if err != nil {
		return err
	}

	fmt.Fprint(rw, string(analysesJSON))
	return nil
}
