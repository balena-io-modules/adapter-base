package main

import (
	"os"

	"github.com/resin-io/adapter-base/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
