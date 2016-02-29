package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/docker"
	//"strings"
	"errors"
)

var pushCmd = &cobra.Command{
	Use:   "push IMAGE",
	Short: "Push Docker image to Clair",
	Long:  `Upload a Docker image to Clair for further analysis`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("hyperclair: \"push\" requires a minimum of 1 argument")
		}

		err := push(args[0])
		if err != nil {
			return err
		}
		fmt.Println("All is ok")
		return nil
	},
}

func push(imageName string) error {
	image, err := docker.Parse(imageName)
	if err != nil {
		return err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := HyperclairURI + "/" + image.Name + "?realm=" + registry + "&reference=" + image.Tag
	response, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}

	return nil
}

func init() {
	RootCmd.AddCommand(pushCmd)
}
