package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/docker"
)

func hyperclairURI() string {
	return viper.GetString("hyperclair.uri") + ":" + strconv.Itoa(viper.GetInt("hyperclair.port")) + "/v1"
}

//Pull call Hyperclair server to get Light Manifest
func Pull(imageName string) (docker.Image, error) {
	image, err := docker.Parse(imageName)
	if err != nil {
		return docker.Image{}, err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := hyperclairURI() + "/" + image.Name + "?realm=" + registry + "&reference=" + image.Tag
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

func Push(imageName string) error {
	image, err := docker.Parse(imageName)
	if err != nil {
		return err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := hyperclairURI() + "/" + image.Name + "?realm=" + registry + "&reference=" + image.Tag
	response, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	return nil
}

func Analyse(imageName string) (clair.ImageAnalysis, error) {
	image, err := docker.Parse(imageName)
	if err != nil {
		return clair.ImageAnalysis{}, err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := hyperclairURI() + "/" + image.Name + "/analysis" + "?realm=" + registry + "&reference=" + image.Tag
	response, err := http.Get(url)
	if err != nil {
		return clair.ImageAnalysis{}, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		if err != nil {
			return clair.ImageAnalysis{}, err
		}
		return clair.ImageAnalysis{}, fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}
	imageAnalysis := clair.ImageAnalysis{}
	json.Unmarshal(body, &imageAnalysis)
	return imageAnalysis, nil
}

func Health() error {

	url := hyperclairURI() + "/health"
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("Hyperclair in Unhealthy state")
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")

	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON.Bytes()))
	return nil
}

func Versions() error {

	url := hyperclairURI() + "/versions"
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")

	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON.Bytes()))
	return nil
}

func Report(imageName string) error {
	analyses, err := Analyse(imageName)

	if err != nil {
		return err
	}
	SaveAnalysisReport(analyses)
	return nil
}
