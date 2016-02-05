package docker

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
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
	log.Println("Pull is Unauthorized")
	return response.StatusCode == 401
}

func Authenticate(dockerResponse *http.Response, request *http.Request) error {
	bearerToken := BearerAuthParams(dockerResponse)
	url := bearerToken["realm"] + "?service=" + bearerToken["service"] + "&scope=" + bearerToken["scope"]
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	setBasicAuth(req)

	response, err := InitClient().Do(req)

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

	setBearerAuthorization(request, tok.String())

	return nil
}

func setBasicAuth(request *http.Request) {
	request.SetBasicAuth(viper.GetString("auth.user"), viper.GetString("auth.password"))
}

func setBearerAuthorization(request *http.Request, token string) {
	request.Header.Set("Authorization", "Bearer "+token)
}

//InitClient create a http.Client with Transport configuration
func InitClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: viper.GetBool("auth.insecureSkipVerify")},
		DisableCompression: true,
	}
	return &http.Client{Transport: tr}
}
