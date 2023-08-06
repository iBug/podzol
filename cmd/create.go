package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/client"
	"github.com/ustclug/podzol/pkg/docker"
	"github.com/ustclug/podzol/pkg/format"
)

var createCmd = &cobra.Command{
	Use:   "create TOKEN APPLICATION IMAGE PORT [timeout]",
	Short: "Create a new container",
	Long:  `Create a new container with the given arguments. If timeout is not specified, it defaults to a minute.`,
	RunE:  createRunE,
}

func createRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 4 || len(args) > 5 {
		return cmd.Help()
	}

	token := args[0]
	userID, err := format.ParseUserID(token)
	if err != nil {
		return err
	}
	application := args[1]
	image := args[2]
	portnum, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	port := uint16(portnum)
	if int(port) != portnum {
		return fmt.Errorf("invalid port number: %d", portnum)
	}

	timeout := time.Minute
	if len(args) == 5 {
		timeout, err = time.ParseDuration(args[4])
		if err != nil {
			return err
		}
	}

	// Command line arguments passed
	cmd.SilenceUsage = true

	c := client.NewClient(viper.GetViper())
	opts := docker.ContainerOptions{
		User:     userID,
		Token:    token,
		AppName:  application,
		Image:    image,
		Port:     port,
		Lifetime: timeout,
	}
	data, err := c.Create(opts)
	if err != nil {
		return err
	}
	return format.ShowContainer(cmd.OutOrStdout(), data)
}

func init() {
	rootCmd.AddCommand(createCmd)
}
