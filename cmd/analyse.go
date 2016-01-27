package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/docker/image"
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

		image, err := image.Parse(args[0])
		if err != nil {
			return err
		}

		if err := image.Pull(); err != nil {
			return err
		}

		fmt.Println("analysing Image")
		if _, err := image.Analyse(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(analyseCmd)
}
