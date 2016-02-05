package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Pull Image from Registry or Hub depending on image name
func (image *Image) Pull() error {
	fmt.Println("Pull image ", image.Name)
	client := InitClient()
	request, err := http.NewRequest("GET", image.ManifestURI(), nil)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if IsUnauthorized(*resp) {
		err := Authenticate(resp, request)

		if err != nil {
			return err
		}

		resp, err = client.Do(request)

		if err != nil {
			return err
		}

	}

	return image.parseManifest(resp)
}

func (image *Image) parseManifest(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}
	err = json.Unmarshal(body, &image)

	if err != nil {
		return err
	}

	image.uniqueLayers()
	return nil
}

func (image *Image) uniqueLayers() {
	encountered := map[Layer]bool{}
	result := []Layer{}

	for index := range image.FsLayers {
		if encountered[image.FsLayers[index]] != true {
			encountered[image.FsLayers[index]] = true
			result = append(result, image.FsLayers[index])
		}
	}
	image.FsLayers = result
}
