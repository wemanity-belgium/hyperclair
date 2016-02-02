package docker

import (
	"fmt"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/utils"
)

func (image *Image) Analyse() ([]clair.Analysis, error) {
	clair.Config()
	layerCount := len(image.FsLayers)
	analysisResult := []clair.Analysis{}

	for index := range image.FsLayers {
		layer := image.FsLayers[layerCount-index-1]

		if analysis, err := clair.AnalyseLayer(layer.BlobSum); err != nil {
			fmt.Printf("Error analysing layer [%v] %d/%d: %v\n", utils.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
		} else {
			fmt.Printf("Analysis [%v] found %d vulnerabilities.\n", utils.Substr(layer.BlobSum, 0, 12), len(analysis.Vulnerabilities))
			analysisResult = append(analysisResult, analysis)
		}
	}
	return analysisResult, nil
}
