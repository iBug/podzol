package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/config"
	"github.com/ustclug/podzol/pkg/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server process",
	Long:  "Run as the server process",
	RunE:  serverRunE,
}

func serverRunE(cmd *cobra.Command, args []string) error {
	if err := config.Load(); err != nil {
		cmd.SilenceUsage = true
		return err
	}

	s, err := server.NewServer(viper.GetViper())
	if err != nil {
		return err
	}
	return s.Run()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
