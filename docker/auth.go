package docker

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
	s := strings.Fields(r.Header.Get("Www-Authenticate"))

	// Adding Basic authentication support
	bearer_or_basic, _ := regexp.MatchString("Bearer|Basic", s[0])
	if len(s) != 2 || !bearer_or_basic {
		logrus.Fatalf("Authentication Realm is not supported: ", s[0])
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
	url := ""

	s := strings.Fields(dockerResponse.Header.Get("Www-Authenticate"))

	if s[0] == "Bearer" {
		url = bearerToken["realm"] + "?service=" + bearerToken["service"]
	} else {
		url = bearerToken["realm"] + "v2/"
	}

	if bearerToken["scope"] != "" {
		url += "&scope=" + bearerToken["scope"]
	}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	l, err := config.GetLogin(strings.Trim(bearerToken["realm"], "/"))
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

	if s[0] == "Bearer" {
		var tok token
		err = json.Unmarshal(body, &tok)

		if err != nil {
			return err
		}
		request.Header.Set("Authorization", "Bearer "+tok.String())
	} else {
		auth := fmt.Sprintf("%s:%s", l.Username, l.Password)
		request.Header.Set("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(auth)))
	}

	return nil
}
