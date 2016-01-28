package image

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/database"
	"github.com/wemanity-belgium/hyperclair/utils"
)

func (im *DockerImage) Push() error {
	if im.Registry != "" {
		return im.pushFromRegistry()
	}

	return im.pushFromHub()
}

func (im *DockerImage) pushFromHub() error {
	return errors.New("Clair Analysis for Docker Hub is not implemented yet!")
}

func (im *DockerImage) pushFromRegistry() error {
	clair.Config()
	layerCount := len(im.Manifest.FsLayers)

	parentID := ""
	for index, layer := range im.Manifest.FsLayers {
		fmt.Printf("Pushing Layer %d/%d\n", index, layerCount)

		database.InsertRegistryMapping(layer.BlobSum, im.Registry)
		payload := clair.Layer{
			ID:          layer.BlobSum,
			Path:        im.BlobsURI(layer.BlobSum),
			ParentID:    parentID,
			ImageFormat: "Docker",
		}
		payload.Path = strings.Replace(payload.Path, im.Registry, "hyperclair:9999", 1)
		log.Println("Path: ", payload.Path)

		if err := clair.AddLayer(payload); err != nil {
			fmt.Printf("Error adding layer [%v] %d/%d: %v\n", utils.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
			parentID = ""
		} else {
			parentID = layer.BlobSum
		}
	}

	return nil
}
