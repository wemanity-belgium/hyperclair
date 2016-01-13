package pull

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jgsqware/hyperclair/services"
)

type layer struct {
	BlobSum string
}

type manifest struct {
	imageName string
	tag       string
	FsLayers  []layer
}

func GetLayers(services services.Services, imageName string, tag string) (manifest, error) {
	url := services.Registry.GetManifestUrl(imageName, tag)

	var manifestObject manifest

	response, err := http.Get(url)
	if err != nil {
		return manifestObject, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return manifestObject, fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &manifestObject)
	//err = json.NewDecoder(strings.NewReader(body)).Decode(&manifestObject)

	if err != nil {
		return manifestObject, err
	}

	uniqueLayers(&manifestObject)
	return manifestObject, nil
}

func uniqueLayers(manifestObject *manifest) {
	encountered := map[layer]bool{}
	result := []layer{}

	for index := range manifestObject.FsLayers {
		if encountered[manifestObject.FsLayers[index]] != true {
			encountered[manifestObject.FsLayers[index]] = true
			result = append(result, manifestObject.FsLayers[index])
		}
	}

	manifestObject.FsLayers = result
}

type token struct {
	Token string
}

func PullForHub() (manifest, error) {
	response, err := http.Get("https://auth.docker.io/token?service=registry.docker.io&scope=repository:gliderlabs/alpine:pull")

	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var tok token
	err = json.Unmarshal(body, &tok)

	client := &http.Client{}

	request, _ := http.NewRequest("GET", "https://registry-1.docker.io/v2/gliderlabs/alpine/manifests/latest", nil)

	request.Header.Add("Authorization", "Bearer "+tok.Token)

	resp, err1 := client.Do(request)

	if err1 != nil {
		fmt.Errorf(err1.Error())
	}

	manifestBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var manifestObject manifest
	err = json.Unmarshal(manifestBody, &manifestObject)

	if err != nil {
		return manifestObject, err
	}

	uniqueLayers(&manifestObject)

	return manifestObject, nil
}
