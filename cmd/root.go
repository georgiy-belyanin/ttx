package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/georgiy-belyanin/ttx/config"
	"github.com/georgiy-belyanin/ttx/runner"
	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "ttx",
	Short: "TTX - Tarantool CLI tool for developers",
	Long: `TTX is a Tarantool CLI tool for developers.

TTX simplifies working with Tarantool configuration and testing the clusters
during the development.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(configPath) == 0 {
			fmt.Fprintf(os.Stderr, "Found no configuration file in the current dir or it's parent directories.\nMake sure there is a configuration file named like \"config.yml\" or supply a path to it using \"-c\" flag.\n")
			return
		}

		err := runner.RunClusterFromConfig(context.Background(), configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to start the cluster based on the configuration path \"%s\": %s\n", configPath, err)
			return
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to a cluster configuration file (default: seek YAML\n\"config\" and \"source\" files in the nearest dirs)")

	if len(configPath) == 0 {
		var err error
		configPath, err = config.FindYamlFileAtPath(".")
		if err != nil {
			configPath = ""
		}
	}
}
