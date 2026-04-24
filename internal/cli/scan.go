package cli

import (
	"netpeek/internal/scanner"
	"time"

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

	p := scanner.NewPrinter(s.Stats.Total)
	p.OutputInitialText(host, ports)
	defer trackTime(time.Now(), "Scanning", p, s)

	_, err = s.Scan(p)

	return err
}

func trackTime(now time.Time, name string, p *scanner.Printer, s *scanner.Scanner) {
	elapsed := time.Since(now).Seconds()
	p.OutputFinalStats(s.Stats.Opened, s.Stats.Closed, elapsed)
}
