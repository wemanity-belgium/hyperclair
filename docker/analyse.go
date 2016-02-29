package docker

import (
	"log"
	"strings"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

//Analyse return Clair Image analysis
func (image *Image) Analyse() clair.ImageAnalysis {
	clair.Config()
	layerCount := len(image.FsLayers)
	analysisResult := []clair.LayerAnalysis{}
	for index := range image.FsLayers {
		layer := image.FsLayers[layerCount-index-1]

		if analysis, err := clair.AnalyseLayer(layer.BlobSum); err != nil {
			log.Printf("Error analysing layer [%v] %d/%d: %v\n", xstrings.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
		} else {
			log.Printf("Analysis [%v] found %d vulnerabilities.\n", xstrings.Substr(layer.BlobSum, 0, 12), len(analysis.Vulnerabilities))
			analysisResult = append(analysisResult, analysis)
		}
	}
	return clair.ImageAnalysis{
		Registry:  strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2"),
		ImageName: image.Name,
		Tag:       image.Tag,
		Layers:    analysisResult,
	}
}
