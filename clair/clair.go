package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

//URI is Clair Uri
var URI string
var Port int
var Link string
var Priority string
var Report ReportConfig

type ReportConfig struct {
	Path   string
	Format string
}

type Layer struct {
	ID, Path, ParentID string
}

type Vulnerability struct {
	ID, Link, Priority, Description, CausedByPackage string
}
type Analysis struct {
	ID              string
	Vulnerabilities []Vulnerability
}

//Config configure Clair from configFile
func Config() {
	URI = viper.GetString("clair.uri")
	Port = viper.GetInt("clair.port")
	Link = viper.GetString("clair.link")
	Priority = viper.GetString("clair.priority")
	Report.Path = viper.GetString("clair.report.path")
	Report.Format = viper.GetString("clair.report.format")
}

func formatURI() string {
	uri := URI
	if Port != 0 {
		uri = ":" + strconv.Itoa(Port)
	}
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}
	if !strings.HasSuffix(uri, "/v1") {
		uri += "/v1"
	}

	return uri
}

func addLayerURI() string {
	return strings.Join([]string{formatURI(), "layers"}, "/")
}

func analyseLayerURI(id string) string {
	return strings.Join([]string{formatURI(), "layers", id, "vulnerabilities"}, "/") + "?minimumPriority=" + Priority
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

	layerJSONPayload, err := json.MarshalIndent(layer, "", "\t")
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", addLayerURI(), bytes.NewBuffer(layerJSONPayload))
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

func AnalyseLayer(id string) (Analysis, error) {

	response, err := http.Get(analyseLayerURI(id))
	if err != nil {
		return Analysis{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return Analysis{}, fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	body, _ := ioutil.ReadAll(response.Body)

	var analysis Analysis

	err = json.Unmarshal(body, &analysis)
	if err != nil {
		return Analysis{}, err
	}
	analysis.ID = id
	return analysis, nil
}
