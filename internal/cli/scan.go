package cli

import (
	"fmt"
	"netpeek/internal/scanner"

	"github.com/spf13/cobra"
)

var ports string
var host string

func init() {
	scanCmd.Flags().StringVarP(&ports, "port", "p", "", "Ports to scan (e.g., 'all', '80', '80-443', '80,82')")
	scanCmd.Flags().StringVarP(&host, "host", "H", "", "Host to scan")

	scanCmd.MarkFlagRequired("host")
	scanCmd.MarkFlagRequired("port")

	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:  "scan",
	RunE: runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	s, err := scanner.NewScanner(host, ports)

	if err != nil {
		return err
	}

	result, err := s.Scan()
	for _, res := range result {
		fmt.Printf("port: %v, opened: %v\n", res.Port, res.Opened)
	}

	return nil
}
