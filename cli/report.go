package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
)

//Report generate Clair Report
func SaveAnalysisReport(analyses clair.ImageAnalysis) error {
	clair.Config()
	imageName := strings.Replace(analyses.ImageName, "/", "-", -1) + "-" + analyses.Tag

	switch clair.Report.Format {
	case "html":
		html, err := analyses.ReportAsHTML()

		if err != nil {
			return err
		}
		return SaveReport(imageName, string(html))
	case "json":
		json, err := analyses.ReportAsJSON()

		if err != nil {
			return err
		}
		return SaveReport(imageName, string(json))
	default:
		return fmt.Errorf("Unsupported Report format: %v", clair.Report.Format)
	}
}

func reportPath() string {
	return viper.GetString("clair.report.path") + "/" + clair.Report.Format
}

func SaveReport(name string, content string) error {
	path := reportPath()
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	reportsName := path + "/analysis-" + name + "." + clair.Report.Format
	f, err := os.Create(reportsName)
	if err != nil {
		return err
	}

	f.WriteString(content)
	fmt.Printf("%v report at %v\n", strings.ToUpper(clair.Report.Format), reportsName)
	return nil
}
