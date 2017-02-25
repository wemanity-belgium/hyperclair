package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/config"
	"github.com/wemanity-belgium/hyperclair/docker/httpclient"
	"github.com/wemanity-belgium/hyperclair/xerrors"
)

type token struct {
	Value string `json:"token"`
}

func (tok token) String() string {
	return tok.Value
}

//BearerAuthParams parse Bearer Token on Www-Authenticate header
func BearerAuthParams(r *http.Response) map[string]string {
	s := strings.SplitN(r.Header.Get("Www-Authenticate"), " ", 2)
	logrus.Debug("auth header is ", r.Header.Get("Www-Authenticate"))
	if len(s) != 2 || s[0] != "Bearer" {
		// other fields might have spaces in them!
		logrus.Warn("couldn't find bearer string")
		return nil
	}
	result := map[string]string{}

	for _, kv := range strings.Split(s[1], ",") {
		parts := strings.Split(kv, "=")
		if len(parts) != 2 {
			continue
		}
		result[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
	}
	return result
}

func AuthenticateResponse(dockerResponse *http.Response, request *http.Request) error {
	bearerToken := BearerAuthParams(dockerResponse)
	logrus.Debug("bearer token is ", bearerToken)
	url := bearerToken["realm"] + "?service=" + url.QueryEscape(bearerToken["service"])
	if bearerToken["scope"] != "" {
		url += "&scope=" + bearerToken["scope"]
	}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}
	l, err := config.GetLogin(request.URL.Host)
	if err != nil {
		return err
	}
	req.SetBasicAuth(l.Username, l.Password)

	response, err := httpclient.Get().Do(req)

	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return xerrors.Unauthorized
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication server response: %v - %v", response.StatusCode, response.Status)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	var tok token
	err = json.Unmarshal(body, &tok)

	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+tok.String())

	return nil
}
