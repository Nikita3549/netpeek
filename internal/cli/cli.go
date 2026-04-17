package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "Coming soon...",
	Args: cobra.NoArgs,
}

func Execute() error {
	return rootCmd.Execute()
}
