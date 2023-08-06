package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/client"
	"github.com/ustclug/podzol/pkg/docker"
	"github.com/ustclug/podzol/pkg/format"
)

var removeCmd = &cobra.Command{
	Use:   "remove { USER | TOKEN } APPLICATION",
	Short: "Remove a container",
	Long:  `Remove a container`,
	RunE:  removeRunE,
}

func removeRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("bad number of arguments")
	}
	user, err := strconv.Atoi(args[0])
	if err != nil {
		token := args[0]
		user, err = format.ParseUserID(token)
		if err != nil {
			return err
		}
	}
	app := args[1]

	// Arguments validated
	cmd.SilenceUsage = true

	opts := docker.ContainerOptions{
		User:    user,
		AppName: app,
	}
	c := client.NewClient(viper.GetViper())
	err = c.Remove(opts)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "OK")
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
