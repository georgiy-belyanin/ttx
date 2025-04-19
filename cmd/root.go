package cmd

import (
	"fmt"

	"github.com/georgiy-belyanin/ttx/runner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ttx",
	Short: "TTX - Tarantool CLI tool for developers",
	Long: `TTX is a Tarantool CLI tool for developers.

TTX simplifies working with Tarantool configuration and testing the clusters
during the development.
`,

	Run: func(cmd *cobra.Command, args []string) {
		err := runner.RunClusterFromNearestConfig()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
