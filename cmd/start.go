package cmd

import (
	"fmt"

	"github.com/georgiy-belyanin/ttx/runner"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <config.yml>",
	Short: "Start Tarantool cluster",
	Long:  `Start Tarantool cluster from the configuration`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := args[0]
		err := runner.RunClusterFromConfig(configPath)
		if err != nil {
			fmt.Println("Unable to start the cluster from the config", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
