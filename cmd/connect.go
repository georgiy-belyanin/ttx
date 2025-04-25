package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/georgiy-belyanin/ttx/config"
	"github.com/georgiy-belyanin/ttx/console"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [user[:password]@]<instance-name or uri>",
	Short: "Connect to Tarantool instance",
	Long: `This command allows one to connect to Tarantool instance via iproto and
control them.

    ttx connect i-001
    ttx connect user:password@i-001
	ttx connect localhost:3301`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "Please, supply strictly one argument in the format of [user[:password]@]<instance-name or uri>\n")
			return
		}

		config, _ := config.LoadYamlFile(configPath)

		err := console.ConnectByString(context.Background(), config, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to instance: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
