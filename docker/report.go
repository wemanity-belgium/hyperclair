package docker

import (
	"fmt"

	"github.com/wemanity-belgium/hyperclair/clair"
)

//Report generate Clair Report
func (image *Image) Report() error {
	clair.Config()
	analysies, err := image.Analyse()
	if err != nil {
		return err
	}
	for _, analysis := range analysies {
		switch clair.Report.Format {
		case "html":
			return analysis.ReportAsHTML()
		case "json":
			return analysis.ReportAsJSON()
		default:
			return fmt.Errorf("Unsupported Report format: %v", clair.Report.Format)
		}
	}
	return nil
}
