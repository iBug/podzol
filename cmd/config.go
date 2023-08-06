package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ustclug/podzol/pkg/config"
)

var writeConfigPath string

var defaultconfigCmd = &cobra.Command{
	Use:   "defaultconfig [-o output]",
	Short: "Generate a default configuration file",
	Long:  `Generate a default configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.Save(writeConfigPath)
	},
}

func init() {
	rootCmd.AddCommand(defaultconfigCmd)

	defaultconfigCmd.Flags().StringVarP(&writeConfigPath, "output", "o", config.ExampleFile, "output file")
}
