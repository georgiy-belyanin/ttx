package cmd

import (
	"context"
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

		ctx := context.Background()

		if len(args) >= 1 {
			configPath := args[0]
			err = runner.RunClusterFromConfig(ctx, configPath)
		} else {
			err = runner.RunClusterFromNearestConfig(ctx)
		}

		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
