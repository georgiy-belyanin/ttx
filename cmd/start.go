package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/georgiy-belyanin/ttx/runner"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [cluster parts...]",
	Short: "Start Tarantool cluster",
	Long: `Start Tarantool cluster from the configuration.

Use arguments to supply cluster parts to start to. It might be instance names,
replicaset names, group names. If no arguments are supplied the command starts
the whole cluster.

Example:

    ttx start replicaset-005 group-003 instance-008`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(configPath) == 0 {
			fmt.Fprintf(os.Stderr, "Found no configuration file in the current dir or it's parent directories.\nMake sure there is a configuration file named like \"config.yml\" or supply a path to it using \"-c\" flag.\n")
			return
		}

		var err error
		ctx := context.Background()

		if len(args) > 0 {
			err = runner.RunClusterPartsFromConfig(ctx, configPath, args)
		} else {
			err = runner.RunClusterFromConfig(ctx, configPath)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to start the cluster based on the configuration path \"%s\": %s\n", configPath, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
