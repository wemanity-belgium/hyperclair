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
var port int
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

//Analysis Clair analysis
type Analysis struct {
	ID              string
	ImageName       string
	Vulnerabilities []Vulnerability
}

//Count vulnarabilities regarding the priority
func (analysis Analysis) Count(priority string) int {
	count := 0
	for _, vulnerability := range analysis.Vulnerabilities {
		if vulnerability.Priority == priority {
			count++
		}
	}

	return count
}

//Config configure Clair from configFile
func Config() {
	uri = viper.GetString("clair.uri")
	port = viper.GetInt("clair.port")
	formatURI()
	priority = viper.GetString("clair.priority")
	Report.Path = viper.GetString("clair.report.path")
	Report.Format = viper.GetString("clair.report.format")
}

func formatURI() string {
	if port != 0 {
		uri += ":" + strconv.Itoa(port)
	}
	if !strings.HasSuffix(uri, "/v1") {
		uri += "/v1"
	}
	return "http://clair:6060/v1"
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
