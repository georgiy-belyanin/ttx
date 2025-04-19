package main

import (
	"os"

	"github.com/georgiy-belyanin/ttx/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(-1)
	}
}
