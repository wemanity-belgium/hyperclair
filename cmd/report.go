package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
)

var reportCmd = &cobra.Command{
	Use:   "report IMAGE",
	Short: "Generate Docker Image vulnerabilities report",
	Long:  `Generate Docker Image vulnerabilities report as HTML or JSON`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("hyperclair: \"report\" requires a minimum of 1 argument")
		}
		if err := report(args[0]); err != nil {
			return err
		}
		return nil
	},
}

func report(imageName string) error {
	analyses, err := Analyse(imageName)
	if err != nil {
		return err
	}
	saveAnalysisReport(analyses)
	return nil
}

func reportPath() string {
	return viper.GetString("clair.report.path") + "/" + clair.Report.Format
}

func saveAnalysisReport(analyses clair.ImageAnalysis) error {
	clair.Config()
	imageName := strings.Replace(analyses.ImageName, "/", "-", -1) + "-" + analyses.Tag
	switch clair.Report.Format {
	case "html":
		html, err := analyses.ReportAsHTML()
		if err != nil {
			return err
		}
		return saveReport(imageName, string(html))
	case "json":
		json, err := analyses.ReportAsJSON()

		if err != nil {
			return err
		}
		return saveReport(imageName, string(json))
	default:
		return fmt.Errorf("Unsupported Report format: %v", clair.Report.Format)
	}
}

func saveReport(name string, content string) error {
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

func init() {
	RootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP("format", "f", "html", "Format for Report [html,json]")
	viper.BindPFlag("clair.report.format", reportCmd.Flags().Lookup("format"))
}
