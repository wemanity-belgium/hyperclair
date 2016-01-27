package image

import (
	"fmt"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/utils"
)

func (im *DockerImage) Analyse() ([]clair.Analysis, error) {
	clair.Config()
	layerCount := len(im.Manifest.FsLayers)
	analysisResult := []clair.Analysis{}

	for index := range im.Manifest.FsLayers {
		layer := im.Manifest.FsLayers[layerCount-index-1]

		if analysis, err := clair.AnalyseLayer(layer.BlobSum); err != nil {
			fmt.Printf("Error analysing layer [%v] %d/%d: %v\n", utils.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
		} else {
			fmt.Printf("Analysis [%v] found %d vulnerabilities.\n", utils.Substr(layer.BlobSum, 0, 12), len(analysis.Vulnerabilities))
			analysis.ImageName = im.GetName()
			analysisResult = append(analysisResult, analysis)
		}
	}
	return analysisResult, nil
}
