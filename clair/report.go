package clair

import (
	"bytes"
	"encoding/json"
	"text/template"
)

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

	templte, err := template.New("analysis-template").ParseFiles("templates/analysis-template.html")
	if err != nil {
		return "", err
	}
	var doc bytes.Buffer
	err = templte.ExecuteTemplate(&doc, "analysis-template.html", analyses)
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
