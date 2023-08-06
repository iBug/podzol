package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ustclug/podzol/pkg/config"
)

var defaultconfigCmd = &cobra.Command{
	Use:   "defaultconfig",
	Short: "Generate a default configuration file",
	Long:  `Generate a default configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.Save(config.ExampleFile)
	},
}

func init() {
	rootCmd.AddCommand(defaultconfigCmd)
}
