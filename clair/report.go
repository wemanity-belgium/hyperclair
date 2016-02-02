package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/wemanity-belgium/hyperclair/utils"
)

//ReportConfig  Reporting configuration
type ReportConfig struct {
	Path   string
	Format string
}

//ReportAsJSON report analysis as Json
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

//ReportAsHTML report analysis as Htll
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
