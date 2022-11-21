package main

import (
	"os"

	"github.com/mdordoy/peloton-to-garmin/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
