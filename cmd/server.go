package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server process",
	Long:  "Run as the server process",
	RunE:  serverRunE,
}

func serverRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
