package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

//URI is Clair Uri
var URI string
var Port int
var Link string

type Layer struct {
	ID, Path, ParentID string
}

//Config configure Clair from configFile
func Config() {
	URI = viper.GetString("clair.uri")
	Port = viper.GetInt("clair.port")
	Link = viper.GetString("clair.link")
}

func formatURI() string {
	uri := URI
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}
	if !strings.HasSuffix(uri, "/v1") {
		uri += "/v1"
	}

	return uri
}

func (layer *Layer) updateLayer() {
	if strings.Contains(layer.Path, "localhost") {
		layer.Path = strings.Replace(layer.Path, "localhost", Link, 1)
	} else if strings.Contains(layer.Path, "127.0.0.1") {
		layer.Path = strings.Replace(layer.Path, "127.0.0.1", Link, 1)
	}
}

func AddLayer(layer Layer) error {
	layer.updateLayer()

	layerJSONPayload, err := json.Marshal(layer)
	if err != nil {
		return err
	}
	fmt.Println(strings.Join([]string{URI, "/layers"}, "/"))
	request, err := http.NewRequest("POST", strings.Join([]string{"http://localhost:6060/v1", "layers"}, "/"), bytes.NewBuffer(layerJSONPayload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 201 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	return nil
}
