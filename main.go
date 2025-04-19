package main

import (
	"fmt"

	"github.com/georgiy-belyanin/ttx/runner"
	"os"
)

func main() {
	configPath := os.Args[1]

	err := runner.RunClusterFromConfig(configPath)
	if err != nil {
		fmt.Println("Unable to start the cluster from the config:", err)
		return
	}
}
