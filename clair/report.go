package clair

import (
	"bytes"
	"encoding/json"
	"text/template"
)

//go:generate go-bindata -pkg clair -o templates.go templates/...

//ReportConfig  Reporting configuration
type ReportConfig struct {
	Path   string
	Format string
}

//ReportAsJSON report analysis as Json
func (analyses ImageAnalysis) ReportAsJSON() ([]byte, error) {

	analysesAsJSON, err := json.MarshalIndent(analyses, "", "\t")
	if err != nil {
		return nil, err
	}
	return analysesAsJSON, nil
}

//ReportAsHTML report analysis as HTML
func (analyses ImageAnalysis) ReportAsHTML() (string, error) {
	asset, err := Asset("templates/analysis-template.html")
	if err != nil {
		return "", err
	}

	templte := template.Must(template.New("analysis-template").Parse(string(asset)))

	var doc bytes.Buffer
	err = templte.Execute(&doc, analyses)
	if err != nil {
		return "", err
	}
	return doc.String(), nil
}

//ReportAsJSON report analysis as Json
func (layerAnalysis LayerAnalysis) ReportAsJSON() (string, error) {
	json, err := json.MarshalIndent(layerAnalysis, "", "\t")
	if err != nil {
		return "", err
	}
	return string(json), nil
}
