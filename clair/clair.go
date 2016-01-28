package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/utils"
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
	ID, Path, ParentID, ImageFormat string
}

type Vulnerability struct {
	ID, Link, Priority, Description, CausedByPackage string
}
type Analysis struct {
	ID              string
	ImageName       string
	Vulnerabilities []Vulnerability
}

func (analysis Analysis) Count(priority string) int {
	count := 0
	for _, vulnerability := range analysis.Vulnerabilities {
		if vulnerability.Priority == priority {
			count++
		}
	}

	return count
}

func (analysis Analysis) ReportAsJSON() error {
	if err := os.MkdirAll("reports/json", 0777); err != nil {
		return err
	}

	reportsName := "reports/json/analysis-" + strings.Replace(utils.Substr(analysis.ID, 0, 12), ":", "", 1) + ".json"
	f, err := os.Create(reportsName)
	if err != nil {
		return err
	}

	defer f.Close()
	json, err := json.MarshalIndent(analysis, "", "\t")
	if err != nil {
		return err
	}
	f.Write(json)
	fmt.Println("JSON report at ", reportsName)
	return nil
}

func (analysis Analysis) ReportAsHTML() error {
	if err := os.MkdirAll("reports/html", 0777); err != nil {
		return err
	}

	t, err := template.New("analysis-template").ParseFiles("templates/analysis-template.html")
	if err != nil {
		return err
	}
	reportsName := "reports/html/analysis-" + strings.Replace(utils.Substr(analysis.ID, 0, 12), ":", "", 1) + ".html"
	f, err := os.Create(reportsName)
	if err != nil {
		return err
	}

	defer f.Close()
	var doc bytes.Buffer
	err = t.ExecuteTemplate(&doc, "analysis-template.html", analysis)
	if err != nil {
		return err
	}
	f.WriteString(doc.String())
	fmt.Println("HTML report at ", reportsName)
	return nil
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
		uri += ":" + strconv.Itoa(Port)
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
