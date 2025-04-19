package cmd

import (
	"fmt"

	"github.com/georgiy-belyanin/ttx/runner"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [config.yml]",
	Short: "Start Tarantool cluster",
	Long:  `Start Tarantool cluster from the configuration`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if len(args) > 1 {
			configPath := args[0]
			err = runner.RunClusterFromConfig(configPath)
		} else {
			err = runner.RunClusterFromNearestConfig()
		}

		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
