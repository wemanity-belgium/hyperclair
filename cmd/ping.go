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
	"github.com/spf13/cobra"
	"log"

	"github.com/jgsqware/hyperclair/services"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the registry",
	Long: `Ping the Docker registry to check if it's up`,
	Run: func(cmd *cobra.Command, args []string) {

		services := services.New()

		//TODO the Get Value is not great
		err := services.Clair.Ping()
		if err != nil {
			log.Printf(err.Error())
		}

		err = services.Registry.Ping()
		if err != nil {
			log.Printf(err.Error())
		}

		fmt.Printf("All is up!")
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
