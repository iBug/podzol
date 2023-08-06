package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ustclug/podzol/pkg"
)

var rootCmd = &cobra.Command{
	Use: strings.ToLower(pkg.Name),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
