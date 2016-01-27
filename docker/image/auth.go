package image

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func (image *DockerImage) login() error {
	loginURI := AuthURI() + "?service=" + image.Registry + "&scope=repository:" + image.GetOnlyName() + ":pull"
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
