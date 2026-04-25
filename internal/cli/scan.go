package cli

import (
	"netpeek/internal/scanner"
	"time"

	"github.com/spf13/cobra"
)

var ports string
var host string
var concurrency int
var timeout int

func init() {
	scanCmd.Flags().StringVarP(&ports, "port", "p", "", "Ports to scan (e.g., 'all', '80', '80-443', '80,82')")
	scanCmd.Flags().StringVarP(&host, "host", "H", "", "Host to scan")
	scanCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 0, "Goroutine workers for scanning")
	scanCmd.Flags().IntVarP(&timeout, "timeout", "t", 0, "Timeout per port scanning")

	scanCmd.MarkFlagRequired("host")
	scanCmd.MarkFlagRequired("port")

	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:  "scan",
	RunE: runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	var options []scanner.Option

	if concurrency > 0 {
		options = append(options, scanner.WithWorkerCount(concurrency))
	}

	if timeout > 0 {
		t := time.Duration(timeout) * time.Millisecond
		options = append(options, scanner.WithDialTimeout(t))
	}

	s, err := scanner.NewScanner(host, ports, options...)
	if err != nil {
		return err
	}

	p := scanner.NewPrinter(s.Stats.Total)
	p.OutputInitialText(host, ports)
	defer trackTime(time.Now(), p, s)

	_, err = s.Scan(p)

	return err
}

func trackTime(now time.Time, p *scanner.Printer, s *scanner.Scanner) {
	elapsed := time.Since(now).Seconds()
	p.OutputFinalStats(s.Stats.Opened, s.Stats.Closed, elapsed)
}
