package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use: "scan",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("scanning...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
