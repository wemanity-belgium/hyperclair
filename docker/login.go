package docker

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/docker/httpclient"
	"github.com/wemanity-belgium/hyperclair/xerrors"
)

//Pull Image from Registry or Hub depending on image name
func Login(registry string) (bool, error) {

	logrus.Info("log in: ", registry)

	client := httpclient.Get()
	request, err := http.NewRequest("GET", registry, nil)
	response, err := client.Do(request)
	if err != nil {
		return false, fmt.Errorf("log in %v: %v", registry, err)
	}
	authorized := response.StatusCode != http.StatusUnauthorized
	if !authorized {
		logrus.Info("Unauthorized access")
		err := AuthenticateResponse(response, request)

		if err != nil {
			if err == xerrors.Unauthorized {
				authorized = false
			}
			return false, err
		} else {
			authorized = true
		}
	}

	return authorized, nil
}
