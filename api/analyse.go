package api

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

func AnalyseHandler(rw http.ResponseWriter, request *http.Request) error {
	local := request.URL.Query()["local"]

	docker.IsLocal = len(local) > 0
	logrus.Debugf("Hyperclair is local: %v", docker.IsLocal)

	var err error
	var image docker.Image

	if !docker.IsLocal {
		image, err = docker.Pull(parseImageURL(request))

		if err != nil {
			return err
		}

	} else {
		image, err = docker.Parse(parseImageURL(request))
		if err != nil {
			return err
		}
		docker.FromHistory(&image)
		if err != nil {
			return err
		}
	}

	analyses := docker.Analyse(image)
	analysesJSON, err := xstrings.ToIndentJSON(analyses)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(analysesJSON))
	return nil
}
