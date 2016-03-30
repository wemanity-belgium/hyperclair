package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

func LocalHandler(rw http.ResponseWriter, request *http.Request) error {
	logrus.Debugln("In Local handler")
	return nil
}
