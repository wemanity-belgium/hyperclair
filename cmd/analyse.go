package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"errors"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var analyseCmd = &cobra.Command{
	Use:   "analyse IMAGE",
	Short: "Analyse Docker image",
	Long:  `Analyse a Docker image with Clair, against Ubuntu, Red hat and Debian vulnerabilities databases`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("hyperclair: \"analyse\" requires a minimum of 1 argument")
		}

		imageAnalysis, err := Analyse(args[0])

		if err != nil {
			return err
		}
		fmt.Printf("Image Analysis:\t %v/%v:%v\n\n", imageAnalysis.Registry, imageAnalysis.ImageName, imageAnalysis.Tag)

		for _, layerAnalysis := range imageAnalysis.Layers {

			fmt.Printf("Analysis [%v] found %d vulnerabilities.\n", xstrings.Substr(layerAnalysis.ID, 0, 12), len(layerAnalysis.Vulnerabilities))
		}
		return nil
	},
}

//Analyse call the clair analysis function and return the Image Analysis
func Analyse(imageName string) (clair.ImageAnalysis, error) {
	image, err := docker.Parse(imageName)
	if err != nil {
		return clair.ImageAnalysis{}, err
	}
	registry := strings.TrimSuffix(strings.TrimPrefix(image.Registry, "http://"), "/v2")
	url := HyperclairURI + "/" + image.Name + "/analysis" + "?realm=" + registry + "&reference=" + image.Tag
	response, err := http.Get(url)
	if err != nil {
		return clair.ImageAnalysis{}, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		if err != nil {
			return clair.ImageAnalysis{}, err
		}
		return clair.ImageAnalysis{}, fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}
	imageAnalysis := clair.ImageAnalysis{}
	json.Unmarshal(body, &imageAnalysis)
	return imageAnalysis, nil
}

func init() {
	RootCmd.AddCommand(analyseCmd)
}
