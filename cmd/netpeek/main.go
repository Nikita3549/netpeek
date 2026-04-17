package main

import (
	"netpeek/internal/cli"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
