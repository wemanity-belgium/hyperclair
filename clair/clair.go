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

var uri string
var link string
var priority string

//Report Reporting Config value
var Report ReportConfig

//Layer Clair Layer
type Layer struct {
	ID, Path, ParentID, ImageFormat string
}

//Vulnerability Clair vulnerabilities
type Vulnerability struct {
	ID, Link, Priority, Description, CausedByPackage string
}

//LayerAnalysis Clair layer analysis
type LayerAnalysis struct {
	ID              string
	Vulnerabilities []Vulnerability
}

//ImageAnalysis Full image analysis
type ImageAnalysis struct {
	Registry  string
	ImageName string
	Tag       string
	Layers    []LayerAnalysis
}

type Health interface{}

func (imageAnalysis ImageAnalysis) Name() string {
	return imageAnalysis.Registry + "/" + imageAnalysis.ImageName + ":" + imageAnalysis.Tag
}

//Count vulnarabilities in all layers regarding the priority
func (imageAnalysis ImageAnalysis) Count(priority string) int {
	count := 0
	for _, layer := range imageAnalysis.Layers {
		count += layer.Count(priority)
	}
	return count
}

//Count vulnarabilities regarding the priority
func (layerAnalysis LayerAnalysis) Count(priority string) int {
	count := 0
	for _, vulnerability := range layerAnalysis.Vulnerabilities {
		if vulnerability.Priority == priority {
			count++
		}
	}

	return count
}

//Config configure Clair from configFile
func Config() {
	formatClairURI()
	priority = viper.GetString("clair.priority")
	Report.Path = viper.GetString("clair.report.path")
	Report.Format = viper.GetString("clair.report.format")
}

func HealthURI() string {
	return uri + "/health"
}
func formatClairURI() {
	uri = viper.GetString("clair.uri")
	port := viper.GetInt("clair.port")

	if port != 0 {
		uri += ":" + strconv.Itoa(port)
	}
	if !strings.HasSuffix(uri, "/v1") {
		uri += "/v1"
	}
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}
}

func addLayerURI() string {
	return strings.Join([]string{uri, "layers"}, "/")
}

func analyseLayerURI(id string) string {
	return strings.Join([]string{uri, "layers", id, "vulnerabilities"}, "/") + "?minimumPriority=" + priority
}

//AddLayer to Clair for analysis
func AddLayer(layer Layer) error {
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

//AnalyseLayer get Analysis os specified layer
func AnalyseLayer(id string) (LayerAnalysis, error) {

	response, err := http.Get(analyseLayerURI(id))
	if err != nil {
		return LayerAnalysis{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return LayerAnalysis{}, fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	body, _ := ioutil.ReadAll(response.Body)

	var analysis LayerAnalysis

	err = json.Unmarshal(body, &analysis)
	if err != nil {
		return LayerAnalysis{}, err
	}
	analysis.ID = id
	return analysis, nil
}

func IsHealthy() (Health, error) {
	Config()
	response, err := http.Get(HealthURI())

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var health Health
	err = json.Unmarshal(body, &health)

	if err != nil {
		return nil, err
	}

	return health, nil
}
