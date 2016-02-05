package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/cli"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Versions of Hyperclair and underlying services",
	Long:  `Get Versions of Hyperclair and underlying services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Versions(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
