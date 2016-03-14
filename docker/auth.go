package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
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
	fmt.Println("Www-Authenticate: ", s)
	fmt.Println("Headers: ", r.Header)
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

func Authenticate(dockerResponse *http.Response, request *http.Request) error {
	bearerToken := BearerAuthParams(dockerResponse)
	url := bearerToken["realm"] + "?service=" + bearerToken["service"] + "&scope=" + bearerToken["scope"]
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	serviceAuthorization := strings.Replace(bearerToken["service"], ".", "_", -1)
	a := viper.Get("auth." + serviceAuthorization)
	if a == nil {
		return fmt.Errorf("no login information for %v", serviceAuthorization)
	}
	authorizations := viper.Sub("auth." + serviceAuthorization)

	user := authorizations.GetString("user")
	password := authorizations.GetString("password")
	req.SetBasicAuth(user, password)

	fmt.Println("req: ", req.URL)
	response, err := httpclient.Get().Do(req)

	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return xerrors.Unauthorized
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
