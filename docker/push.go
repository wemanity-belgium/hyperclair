package docker

import (
	"fmt"
	"strings"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/database"
	"github.com/wemanity-belgium/hyperclair/str"
)

//Push image to Clair for analysis
func (image Image) Push() error {
	clair.Config()
	layerCount := len(image.FsLayers)

	parentID := ""
	for index, layer := range image.FsLayers {
		fmt.Printf("Pushing Layer %d/%d\n", index, layerCount)

		database.InsertRegistryMapping(layer.BlobSum, image.Registry)
		payload := clair.Layer{
			ID:          layer.BlobSum,
			Path:        image.BlobsURI(layer.BlobSum),
			ParentID:    parentID,
			ImageFormat: "Docker",
		}
		//FIXME Update to TLS
		payload.Path = strings.Replace(payload.Path, image.Registry, "http://hyperclair:9999/v2", 1)
		if err := clair.AddLayer(payload); err != nil {
			fmt.Printf("Error adding layer [%v] %d/%d: %v\n", str.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
			parentID = ""
		} else {
			parentID = layer.BlobSum
		}
	}

	return nil
}
