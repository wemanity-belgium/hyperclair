package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/cli"
	"github.com/wemanity-belgium/hyperclair/utils"
	//"strings"
	"errors"
)

var analyseCmd = &cobra.Command{
	Use:   "analyse IMAGE",
	Short: "Analyse Docker image",
	Long:  `Analyse a Docker image with Clair, against Ubuntu, Red hat and Debian vulnerabilities databases`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("hyperclair: \"analyse\" requires a minimum of 1 argument")
		}

		imageAnalysis, err := cli.Analyse(args[0])

		if err != nil {
			return err
		}
		fmt.Printf("Image Analysis:\t %v/%v:%v\n\n", imageAnalysis.Registry, imageAnalysis.ImageName, imageAnalysis.Tag)

		for _, layerAnalysis := range imageAnalysis.Layers {

			fmt.Printf("Analysis [%v] found %d vulnerabilities.\n", utils.Substr(layerAnalysis.ID, 0, 12), len(layerAnalysis.Vulnerabilities))
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(analyseCmd)
}
