package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/client"
	"github.com/ustclug/podzol/pkg/docker"
	"github.com/ustclug/podzol/pkg/format"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List containers",
	Long:  "List containers",
	RunE:  listRunE,

	SilenceUsage: true,
}

func listRunE(cmd *cobra.Command, args []string) error {
	c := client.NewClient(viper.GetViper())

	var opts docker.ContainerOptions
	data, err := c.List(opts)
	if err != nil {
		return err
	}
	return format.ListContainers(cmd.OutOrStdout(), data)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
