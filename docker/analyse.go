package docker

import (
	"log"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

//Analyse return Clair Image analysis
func Analyse(image Image) clair.ImageAnalysis {
	c := len(image.FsLayers)
	res := []clair.LayerAnalysis{}

	for i := range image.FsLayers {
		l := image.FsLayers[c-i-1].BlobSum
		lShort := xstrings.Substr(l, 0, 12)

		if a, err := clair.Analyse(l); err != nil {
			log.Printf("analysing layer [%v] %d/%d: %v", lShort, i+1, c, err)
		} else {
			log.Printf("analysing layer [%v] %d/%d", lShort, i+1, c)
			res = append(res, a)
		}
	}
	return clair.ImageAnalysis{
		Registry:  xstrings.TrimPrefixSuffix(image.Registry, "http://", "/v2"),
		ImageName: image.Name,
		Tag:       image.Tag,
		Layers:    res,
	}
}
