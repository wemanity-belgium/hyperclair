package docker

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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

	if len(s) != 2 || s[0] != "Bearer" {
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

//IsUnauthorized check if the StatusCode is 401
func IsUnauthorized(response http.Response) bool {
	return response.StatusCode == 401
}

func authenticate(dockerResponse *http.Response, request *http.Request) error {
	bearerToken := BearerAuthParams(dockerResponse)

	req, err := http.NewRequest("GET", bearerToken["realm"]+"?service="+bearerToken["service"]+"&scope="+bearerToken["scope"], nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth("jgsqware", string("jgsqware"))

	response, err := initClient().Do(req)

	if err != nil {
		return err
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

func initClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	return &http.Client{Transport: tr}
}
