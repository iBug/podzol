package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/config"
	"github.com/ustclug/podzol/pkg/server"
)

var serverCmd = &cobra.Command{
	Use:   "server [-l listen]",
	Short: "Run server process",
	Long:  "Run as the server process",
	RunE:  serverRunE,
}

func serverRunE(cmd *cobra.Command, args []string) error {
	if err := config.Load(); err != nil {
		cmd.SilenceUsage = true
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Fprintf(cmd.ErrOrStderr(), "Use `%s defaultconfig` to generate a default config.\n", cmd.Root().Name())
		}
		return err
	}

	s, err := server.NewServer(viper.GetViper())
	if err != nil {
		return err
	}

	err = s.DockerInit(cmd.Context())
	if err != nil {
		return err
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Run()
	}()
	go func() {
		errCh <- s.RunHTTP()
	}()
	return <-errCh
}

func init() {
	rootCmd.AddCommand(serverCmd)

	flags := serverCmd.Flags()
	flags.StringP("listen", "l", "", "override listen address")
	viper.BindPFlag("listen-addr", flags.Lookup("listen"))
}
