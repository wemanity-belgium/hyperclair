package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from a Docker registry",
	Long:  `Log out from a Docker registry`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 1 {
			fmt.Println("Only one argument is allowed")
			os.Exit(1)
		}
		var reg string = docker.DockerHub

		if len(args) == 1 {
			reg = args[0]
		}

		if _, err := os.Stat(hyperclairHome()); err == nil {
			f, err := ioutil.ReadFile(hyperclairHome())
			if err != nil {
				fmt.Println(xerrors.InternalError)
				logrus.Fatalf("reading hyperclair file: %v", err)
			}

			var users userMapping
			json.Unmarshal(f, &users)

			delete(users, reg)
			s, err := xstrings.ToIndentJSON(users)
			if err != nil {
				fmt.Println(xerrors.InternalError)
				logrus.Fatalf("indenting logout: %v", err)
			}
			ioutil.WriteFile(hyperclairHome(), s, os.ModePerm)
			fmt.Println("Log out successful")
		} else {
			fmt.Println("You are not logged in")
		}
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
