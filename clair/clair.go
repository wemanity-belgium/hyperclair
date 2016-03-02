package clair

import (
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var uri string
var priority string

//Report Reporting Config value
var Report ReportConfig

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

func (imageAnalysis ImageAnalysis) String() string {
	return imageAnalysis.Registry + "/" + imageAnalysis.ImageName + ":" + imageAnalysis.Tag
}

//Count vulnarabilities in all layers regarding the priority
func (imageAnalysis ImageAnalysis) Count(priority string) int {
	var count int
	for _, layer := range imageAnalysis.Layers {
		count += layer.Count(priority)
	}
	return count
}

//Count vulnarabilities regarding the priority
func (layerAnalysis LayerAnalysis) Count(priority string) int {
	var count int
	for _, vulnerability := range layerAnalysis.Vulnerabilities {
		if vulnerability.Priority == priority {
			count++
		}
	}

	return count
}

func (l LayerAnalysis) ShortName() string {
	return xstrings.Substr(l.ID, 0, 12)
}

func fmtURI(u string, port int) {
	uri = u
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

//Config configure Clair from configFile
func Config() {
	fmtURI(viper.GetString("clair.uri"), viper.GetInt("clair.port"))
	priority = viper.GetString("clair.priority")
	Report.Path = viper.GetString("clair.report.path")
	Report.Format = viper.GetString("clair.report.format")
}
