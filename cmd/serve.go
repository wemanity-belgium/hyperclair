package cmd

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/api/server"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Create hyperclair Server",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		sURL := fmt.Sprintf(":%d", viper.GetInt("hyperclair.port"))
		if local {
			sURL = fmt.Sprintf(":%d", 60000)
		}
		err := server.ListenAndServe(sURL)

		return err
	},
}

func getHyperclairURI(imageName string, path ...string) (string, error) {
	image, err := docker.Parse(imageName)
	if err != nil {
		return "", err
	}
	registry := xstrings.TrimPrefixSuffix(image.Registry, "http://", "/v2")
	registry = xstrings.TrimPrefixSuffix(registry, "https://", "/v2")
	url := fmt.Sprintf("%v/%v", HyperclairURI, image.Name)
	for _, p := range path {
		url = fmt.Sprintf("%v/%v", url, p)
	}

	if local {
		registry = "localhost:60000"
	}

	url = fmt.Sprintf("%v?realm=%v&reference=%v", url, registry, image.Tag)

	return url, nil
}

//StartLocalServer start the hyperclair local server needed for reverse proxy
func StartLocalServer() {
	var err error
	if err != nil {
		fmt.Println(xerrors.InternalError)
		logrus.Fatalf("retrieving internal server URI: %v", err)
	}
	sURL, err := docker.LocalServerIP()
	if err != nil {
		fmt.Println(xerrors.InternalError)
		logrus.Fatalf("retrieving internal server IP: %v", err)
	}
	HyperclairURI = "http://" + sURL + "/v1"
	if err != nil {
		fmt.Println(xerrors.InternalError)
		logrus.Fatalf("starting local server: %v", err)
	}
	err = server.Serve(sURL)
	if err != nil {
		fmt.Println(xerrors.InternalError)
		logrus.Fatalf("starting local server: %v", err)
	}
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolVarP(&local, "local", "l", false, "Use local images")

}
