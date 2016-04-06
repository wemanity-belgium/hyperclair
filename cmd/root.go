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

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
)

var cfgFile string
var logLevel string

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
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level [Panic,Fatal,Error,Warn,Info,Debug]")
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
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("hyperclair: config file not found")
		os.Exit(1)
	}
	logrus.Debugf("Using config file: %v", viper.ConfigFileUsed())

	if viper.Get("clair.uri") == nil {
		viper.Set("clair.uri", "http://localhost")
	}
	if viper.Get("clair.port") == nil {
		viper.Set("clair.port", "6060")
	}
	if viper.Get("clair.healthPort") == nil {
		viper.Set("clair.healthPort", "6061")
	}
	if viper.Get("clair.priority") == nil {
		viper.Set("clair.priority", "Low")
	}
	if viper.Get("clair.report.path") == nil {
		viper.Set("clair.report.path", "reports")
	}
	if viper.Get("clair.report.format") == nil {
		viper.Set("clair.report.format", "html")
	}
	if viper.Get("auth.insecureSkipVerify") == nil {
		viper.Set("auth.insecureSkipVerify", "true")
	}
	if viper.Get("hyperclair.ip") == nil {
		viper.Set("hyperclair.ip", "")
	}
	if viper.Get("hyperclair.port") == nil {
		viper.Set("hyperclair.port", 60000)
	}
	if viper.Get("hyperclair.tempFolder") == nil {
		viper.Set("hyperclair.tempFolder", "/tmp/hyperclair")
	}

	lvl := logrus.WarnLevel
	if logLevel != "" {
		var err error
		lvl, err = logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Warningf("Wrong Log level %v, defaults to [Warning]", logLevel)
			lvl = logrus.WarnLevel
		}
	}
	logrus.SetLevel(lvl)
	clair.Config()
}
