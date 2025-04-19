package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ttx",
	Short: "TTX - Tarantool CLI tool for developers",
	Long: `TTX is a Tarantool CLI tool for developers.

TTX simplifies working with Tarantool configuration and testing the clusters
during the development.
`,
}

func Execute() error {
	return rootCmd.Execute()
}
