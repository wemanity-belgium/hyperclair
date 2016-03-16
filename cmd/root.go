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
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
)

var cfgFile string

//HyperclairURI is the hyperclair server URI. As <hyperclair.uri>:<hypeclair.port>/v1
var HyperclairURI string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hyperclair",
	Short: "Analyse your docker image with Clair, directly from your registry.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.hyperclair.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetEnvPrefix("hyperclair")
	viper.SetConfigName(".hyperclair") // name of config file (without extension)
	viper.AddConfigPath(".")           // adding home directory as first search path
	viper.AutomaticEnv()               // read in environment variables that match
	if cfgFile == "" {                 // enable ability to specify config file via flag
		cfgFile = viper.GetString("config")
	}
	viper.SetConfigFile(cfgFile)
	viper.SetDefault("clair.uri", "http://localhost")
	viper.SetDefault("clair.port", "6060")
	viper.SetDefault("clair.healthPort", "6061")
	viper.SetDefault("clair.priority", "Low")
	viper.SetDefault("clair.report.path", "reports")
	viper.SetDefault("clair.report.format", "html")

	viper.SetDefault("auth.insecureSkipVerify", "false")

	viper.SetDefault("hyperclair.uri", "http://localhost")
	viper.SetDefault("hyperclair.port", "9999")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("hyperclair: config file not found")
		os.Exit(1)
	}
	glog.Info("Using config file:", viper.ConfigFileUsed())
	clair.Config()

	HyperclairURI = viper.GetString("hyperclair.uri") + ":" + strconv.Itoa(viper.GetInt("hyperclair.port")) + "/v1"
}
