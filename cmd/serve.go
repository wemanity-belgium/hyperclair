package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/api/server"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Create hyperclair Server",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := server.ListenAndServe()

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

	url = fmt.Sprintf("%v?realm=%v&reference=%v", url, registry, image.Tag)

	return url, nil
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
