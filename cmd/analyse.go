package cmd

import (
	"fmt"

	"github.com/jgsqware/hyperclair/docker/image"
	"github.com/spf13/cobra"
	//"strings"
	"errors"
)

var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "analyse images",
	Long:  `analyse a Docker image to Clair`,
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
		if err := image.Analyse(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(analyseCmd)
}
