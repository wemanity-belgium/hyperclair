package image

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func (image *DockerImage) login() error {
	loginURI := AuthURI() + "?service=docker_registry&scope=repository:" + image.GetOnlyName() + ":pull"
	// loginURI := AuthURI() + "?service=" + image.Registry + "&scope=repository:" + image.GetOnlyName() + ":pull"
	fmt.Println("LoginURI: ", loginURI)
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter login: ")

	reader.Scan()
	login := reader.Text()
	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(0)

	client := &http.Client{}

	// im.Registry = "https://registry-1.docker.io"

	request, err := http.NewRequest("GET", loginURI, nil)

	if err != nil {
		return err
	}

	request.SetBasicAuth(login, string(password))

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(body, &image.Token)

	return err

}

//BearerAuthParams parse Bearer Token on Www-Authenticate header
func BearerAuthParams(r *http.Response) map[string]string {
	s := strings.Fields(r.Header.Get("Www-Authenticate"))

	if len(s) != 2 || s[0] != "Bearer" {
		return nil
	}
	result := map[string]string{}

	for _, kv := range strings.Split(s[1], ",") {
		fmt.Println("split: ", kv)
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
