package cmd

import (
	"fmt"

	"github.com/jgsqware/hyperclair/docker/image"
	"github.com/spf13/cobra"
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

		image, err := image.Parse(args[0])
		if err != nil {
			return err
		}

		if err := image.Pull(); err != nil {
			return err
		}

		fmt.Println("Pushing Image")
		if err := image.Push(); err != nil {
			return err
		}

		fmt.Println("All is ok")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
}
