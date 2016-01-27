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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/wemanity-belgium/hyperclair/pull"
	// "github.com/wemanity-belgium/hyperclair/config"
	"github.com/wemanity-belgium/hyperclair/services"
	"github.com/wemanity-belgium/hyperclair/utils"

	"errors"
)

// pingCmd represents the ping command
var pullCmd = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull images",
	Long:  `Pull a Docker image`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manifests, err := pull.PullForHub()
		for _, layer := range manifests.FsLayers {
			fmt.Printf("Layer: %s\n", layer.BlobSum)
		}

		os.Exit(0)
		//TODO how to use args with viper
		if len(args) != 1 {
			return errors.New("hyperclair: \"pull\" requires a minimum of 1 argument")
		}

		imageName, tag := utils.SplitImageName(args[0])
		services := services.New()

		manifest, err := pull.GetLayers(services, imageName, tag)
		if err != nil {
			log.Printf(err.Error())
			return err
		}

		for _, layer := range manifest.FsLayers {
			fmt.Printf("Layer: %s\n", layer.BlobSum)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(pullCmd)
}
