// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/docker"
)

// pingCmd represents the ping command
var pullCmd = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull Docker image information",
	Long:  `Pull image information from Docker Hub or Registry`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//TODO how to use args with viper
		if len(args) != 1 {
			return errors.New("hyperclair: \"pull\" requires a minimum of 1 argument")
		}

		image, err := pull(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Layers found: %d\n", len(image.FsLayers))
		for _, layer := range image.FsLayers {
			fmt.Printf("Layer: %s\n", layer.BlobSum)
		}

		return nil
	},
}

func pull(imageName string) (docker.Image, error) {
	image, err := docker.Parse(imageName)
	if err != nil {
		return docker.Image{}, err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := HyperclairURI + "/" + image.Name + "?realm=" + registry + "&reference=" + image.Tag
	response, err := http.Get(url)

	if err != nil {
		return docker.Image{}, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return docker.Image{}, err
	}

	err = json.Unmarshal(body, &image)

	if err != nil {
		return docker.Image{}, err
	}

	return image, nil
}

func init() {
	RootCmd.AddCommand(pullCmd)
}
