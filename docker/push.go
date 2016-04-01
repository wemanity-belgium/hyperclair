package docker

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
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

	if layerCount == 0 {
		logrus.Warningln("there is no layer to push")
	}

	hURL := fmt.Sprintf("http://hyperclair:%d/v2", viper.GetInt("hyperclair.port"))
	if IsLocal {
		var err error
		hURL, err = LocalServerIP()
		hURL = "http://" + hURL + "/v1/local"
		if err != nil {
			return err
		}
		logrus.Infof("using %v as local url", hURL)
	}

	for index, layer := range image.FsLayers {
		lUID := xstrings.Substr(layer.BlobSum, 0, 12)
		logrus.Infof("Pushing Layer %d/%d [%v]", index+1, layerCount, lUID)

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
		if IsLocal {
			payload.Layer.Name = layer.History
			payload.Layer.Path += "/layer.tar"
		}
		payload.Layer.Path = strings.Replace(payload.Layer.Path, image.Registry, hURL, 1)
		if err := clair.Push(payload); err != nil {
			logrus.Infof("adding layer %d/%d [%v]: %v", index+1, layerCount, lUID, err)
			if err != clair.OSNotSupported {
				return err
			}
			parentID = ""
		} else {
			parentID = payload.Layer.Name
		}
	}
	cleanLocal()
	return nil
}
