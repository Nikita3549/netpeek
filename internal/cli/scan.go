package cli

import (
	"fmt"
	"netpeek/internal/scanner"
	"time"

	"github.com/spf13/cobra"
)

var ports string
var host string
var concurrency int
var timeout int

func init() {
	scanCmd.Flags().StringVarP(&ports, "port", "p", "", "Ports to scan: 'all', '80', '1-1024', '80,443,8080'")
	scanCmd.Flags().StringVarP(&host, "host", "H", "", "Host to scan")
	scanCmd.Flags().IntVarP(&concurrency, "concurrency", "c", scanner.DefaultWorkerCount, fmt.Sprintf("Number of concurrent goroutines (default %v)", scanner.DefaultWorkerCount))
	defaultDialTimeout := int(scanner.DefaultDialTimeout.Milliseconds())
	scanCmd.Flags().IntVarP(&timeout, "timeout", "t", defaultDialTimeout, fmt.Sprintf("Timeout per port in milliseconds (default %v)", defaultDialTimeout))

	scanCmd.MarkFlagRequired("host")
	scanCmd.MarkFlagRequired("port")

	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan ports on a target host",
	Long:  "Scan TCP ports on a given host. Supports single ports, ranges, and comma-separated lists.",
	Example: `  netpeek scan -H localhost -p all
  netpeek scan -H 192.168.1.1 -p 80,443
  netpeek scan -H example.com -p 1-1024 -c 500 -t 1000`,
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
