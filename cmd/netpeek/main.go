package main

import (
	"os"

	"github.com/Nikita3549/netpeek/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
