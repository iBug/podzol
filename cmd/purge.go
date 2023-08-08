package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/client"
	"github.com/ustclug/podzol/pkg/format"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge expired containers",
	Long:  `Purge expired containers.`,
	RunE:  purgeRunE,
	Args:  cobra.NoArgs,
}

func purgeRunE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	c := client.NewClient(viper.GetViper())
	infos, err := c.Purge()

	w := cmd.OutOrStdout()
	format.ListContainers(w, infos)
	if err != nil {
		format.ListContainerActionErrors(w, err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}
