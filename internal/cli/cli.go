package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netpeek",
	Short: "Fast TCP port scanner",
	Long:  "netpeek is a fast and lightweight TCP port scanner written in Go.",
	Args:  cobra.NoArgs,
}

func Execute() error {
	return rootCmd.Execute()
}
