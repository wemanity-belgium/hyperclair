package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/cli"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get Health of Hyperclair and underlying services",
	Long:  `Get Health of Hyperclair and underlying services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Health(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(healthCmd)
}
