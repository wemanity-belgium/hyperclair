package image

import (
	"fmt"

	"github.com/wemanity-belgium/hyperclair/clair"
)

func (im *DockerImage) Report() error {
	clair.Config()
	analysies, err := im.Analyse()
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
