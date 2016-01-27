package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jgsqware/hyperclair/utils"
)

func (image DockerImage) isReachable() (int, error) {
	statusCode, err := utils.Ping(formatURI(image.Registry))

	if err != nil {
		return statusCode, errors.New("Registry is not reachable: " + err.Error())
	}

	return statusCode, nil
}

//Pull Image from Registry or Hub depending on image name
func (image *DockerImage) Pull() error {
	if image.Registry != "" {
		return image.pullFromRegistry()
	}

	return image.pullFromHub()
}

func (image *DockerImage) pullFromRegistry() error {
	fmt.Println("Pull image from Registry")

	statusCode, err := image.isReachable()

	if err != nil {
		return err
	}

	client := &http.Client{}
	request, _ := http.NewRequest("GET", image.ManifestURI(), nil)

	if statusCode == 401 {
		image.login()
		request.Header.Add("Authorization", "Bearer "+image.Token.Token)
	}

	response, err := client.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return errors.New(image.GetName() + " not found")
	}

	return image.parseManifest(response)
}

func (image *DockerImage) pullFromHub() error {
	fmt.Println("Pull image from Hub")
	response, err := http.Get(image.AuthURI())

	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var tok token
	err = json.Unmarshal(body, &tok)

	client := &http.Client{}

	image.Registry = "https://registry-1.docker.io"

	request, _ := http.NewRequest("GET", image.ManifestURI(), nil)

	request.Header.Add("Authorization", "Bearer "+tok.Token)

	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	return image.parseManifest(resp)
}

func (image *DockerImage) parseManifest(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}
	err = json.Unmarshal(body, &image.Manifest)

	if err != nil {
		return err
	}

	image.Manifest.uniqueLayers()
	return nil
}

func (manifestObject *DockerManifest) uniqueLayers() {
	encountered := map[Layer]bool{}
	result := []Layer{}

	for index := range manifestObject.FsLayers {
		if encountered[manifestObject.FsLayers[index]] != true {
			encountered[manifestObject.FsLayers[index]] = true
			result = append(result, manifestObject.FsLayers[index])
		}
	}
	manifestObject.FsLayers = result
}
