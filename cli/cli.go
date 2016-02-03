package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/docker"
)

//Pull call Hyperclair server to get Light Manifest
func Pull(imageName string) (docker.Image, error) {
	image, err := docker.Parse(imageName)
	if err != nil {
		return docker.Image{}, err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := hyperclairURI() + "/v1/" + image.Name + "?realm=" + registry + "&reference=" + image.Tag
	response, err := http.Get(url)

	if err != nil {
		return docker.Image{}, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return docker.Image{}, err
	}

	err = json.Unmarshal(body, &image)

	if err != nil {
		return docker.Image{}, err
	}

	return image, nil
}

func hyperclairURI() string {
	return viper.GetString("hyperclair.uri") + ":" + strconv.Itoa(viper.GetInt("hyperclair.port"))
}
