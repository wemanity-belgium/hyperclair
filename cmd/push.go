package cmd

import (
  "github.com/spf13/cobra"
  "github.com/jgsqware/hyperclair/utils"
  "github.com/jgsqware/hyperclair/services"
  //"strings"
  "errors"
  "fmt"
)

var pushCmd = &cobra.Command{
	Use:   "push IMAGE",
	Short: "Push images",
	Long: `Push a Docker image to Clair`,
	RunE: func(cmd *cobra.Command, args []string) error{

		//TODO how to use args with viper
		if len(args) != 1 {
			return errors.New("hyperclair: \"push\" requires a minimum of 1 argument.")
		}

    services := services.New()

    println("Clair url:" + services.Clair.GetUrl())
    println("Registry url:" + services.Registry.GetUrl())

    imageName, tag := utils.SplitImageName(args[0])
    fmt.Printf("Imagename: %s && tag: %s", imageName, tag)


		return nil
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
}
