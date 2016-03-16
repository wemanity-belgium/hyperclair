package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/coreos/clair/api/v1"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/database"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

//Push image to Clair for analysis
func Push(image Image) error {
	layerCount := len(image.FsLayers)

	parentID := ""
	for index, layer := range image.FsLayers {
		lUID := xstrings.Substr(layer.BlobSum, 0, 12)
		fmt.Printf("Pushing Layer %d/%d [%v]\n", index+1, layerCount, lUID)

		err := database.InsertRegistryMapping(layer.BlobSum, image.Registry)
		if err != nil {
			return err
		}
		payload := v1.LayerEnvelope{Layer: &v1.Layer{
			Name:       layer.BlobSum,
			Path:       image.BlobsURI(layer.BlobSum),
			ParentName: parentID,
			Format:     "Docker",
		}}
		//FIXME Update to TLS
		hURL := fmt.Sprintf("http://hyperclair:%d/v2", viper.GetInt("hyperclair.port"))
		payload.Layer.Path = strings.Replace(payload.Layer.Path, image.Registry, hURL, 1)
		if err := clair.Push(payload); err != nil {
			log.Printf("adding layer %d/%d [%v]: %v\n", index+1, layerCount, lUID, err)
			if err != clair.OSNotSupported {
				return err
			}
			parentID = ""
		} else {
			parentID = layer.BlobSum
		}
	}

	return nil
}
