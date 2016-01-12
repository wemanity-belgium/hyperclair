package cmd

import (
	"log"

	"github.com/jgsqware/hyperclair/docker/image"
	"github.com/spf13/cobra"
	//"strings"
	"errors"
	"fmt"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push images",
	Long:  `Push a Docker image to Clair`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("hyperclair: \"push\" requires a minimum of 1 argument")
		}

		image, err := image.Parse(args[0])

		if err != nil {
			log.Printf(err.Error())
			return err
		}

		fmt.Printf("Docker Image: %s", image.GetName())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
}
