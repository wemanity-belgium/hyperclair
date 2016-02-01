package image

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/utils"
)

func (image DockerImage) isReachable() (int, error) {
	statusCode, err := utils.Ping(image.RegistryURI())

	if err != nil {
		return statusCode, errors.New("Registry is not reachable: " + err.Error())
	}

	return statusCode, nil
}

//Pull Image from Registry or Hub depending on image name
func (image *DockerImage) Pull() error {
	fmt.Println("Pull image: ", image.String())

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("http://registry:5000/v2/jgsqware/ubuntu-git/manifests/latest")

	if err != nil {
		return err
	}

	if IsUnauthorized(*resp) {
		bearerToken := BearerAuthParams(resp)

		request, err := http.NewRequest("GET", bearerToken["realm"]+"?service="+bearerToken["service"]+"&scope="+bearerToken["scope"], nil)

		if err != nil {
			return err
		}

		request.SetBasicAuth("jgsqware", string("jgsqware"))

		response, err := client.Do(request)

		if err != nil {
			return err
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &image.Token)

		if err != nil {
			return err
		}

		request, err = http.NewRequest("GET", "http://registry:5000/v2/jgsqware/ubuntu-git/manifests/latest", nil)
		request.Header.Add("Authorization", "Bearer "+image.Token.Token)

		resp, err = client.Do(request)

		if err != nil {
			return err
		}

	}

	return image.parseManifest(resp)
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
	fmt.Println("Manifest: ", string(body))
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
