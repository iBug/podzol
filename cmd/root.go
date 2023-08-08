package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg"
)

var overrideConfigFile string

var rootCmd = &cobra.Command{
	Use:     strings.ToLower(pkg.Name),
	Version: pkg.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if overrideConfigFile != "" {
			viper.SetConfigFile(overrideConfigFile)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&overrideConfigFile, "config", "c", "", "override config file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}
