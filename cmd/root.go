package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ustclug/podzol/pkg"
)

var rootCmd = &cobra.Command{
	Use: pkg.Name,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
