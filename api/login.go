package api

import (
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xerrors"
)

func LoginHandler(rw http.ResponseWriter, request *http.Request) error {
	reg := request.URL.Query()["realm"]

	if len(reg) == 0 {
		return fmt.Errorf("no registry provided")
	}

	logged, err := docker.Login(reg[0])

	if err != nil {
		return err
	}

	if !logged {
		return xerrors.Unauthorized
	}
	return nil
}
