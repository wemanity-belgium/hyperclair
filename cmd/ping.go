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
	"github.com/spf13/viper"

	"log"

	"github.com/jgsqware/hyperclair/ping"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the registry",
	Long: `Ping the Docker registry to check if it's up`,
	Run: func(cmd *cobra.Command, args []string) {

		services := ping.Services{
			ClairURI : viper.GetString("clair.uri"),
			ClairPort : viper.GetInt("clair.port"),
			RegistryURI : viper.GetString("registry.uri"),
			RegistryPort : viper.GetInt("registry.port"),
		}

		//TODO the Get Value is not great
		err := ping.Clair(services.ClairURI,services.ClairPort)
		if err != nil {
			log.Printf(err.Error())
		}

		err = ping.Registry(services.RegistryURI,services.RegistryPort)
		if err != nil {
			log.Printf(err.Error())
		}

		fmt.Printf("All is up!")
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
