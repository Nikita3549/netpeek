package scanner

import (
	"cmp"
	"fmt"
	"net"
	"net/netip"
	"slices"
	"sync"
	"time"
)

const defaultDialTimeout = 500 * time.Millisecond
const defaultWorkerCount = 1024

type Scanner struct {
	Stats   Stats
	address netip.Addr
	ports   []PortRange
	config  Config
}

type Stats struct {
	Total  uint16
	Opened uint16
	Closed uint16
}

type ScannedPort struct {
	Port   uint16
	Opened bool
}

func NewScanner(host, ports string, options ...Option) (*Scanner, error) {
	ip, err := resolveIP(host)
	if err != nil {
		return nil, err
	}

	portRanges, totalPorts, err := parsePorts(ports)
	if err != nil {
		return nil, err
	}

	conf := Config{workerCount: defaultWorkerCount, dialTimeout: defaultDialTimeout}
	for _, option := range options {
		option(&conf)
	}

	return &Scanner{
		ports:   portRanges,
		address: ip,
		config:  conf,
		Stats: Stats{
			Total:  uint16(totalPorts),
			Closed: 0,
			Opened: 0,
		},
	}, nil
}

func (s *Scanner) Scan(p *Printer) ([]ScannedPort, error) {
	jobs := make(chan uint16, s.config.workerCount)
	results := make(chan ScannedPort)
	scannedPorts := make([]ScannedPort, 0)
	var wg sync.WaitGroup

	// Worker
	for i := 0; i < s.config.workerCount; i++ {
		wg.Add(1)
		go s.worker(jobs, results, &wg)
	}

	// Collector
	done := make(chan bool)
	go func() {
		for result := range results {
			scannedPorts = append(scannedPorts, result)
			s.increaseStats(result.Opened)
			p.OutputPort(result.Port, result.Opened)
		}
		done <- true
	}()

	// Producer
	for _, portRange := range s.ports {
		for port := int(portRange.Start); port <= int(portRange.End); port++ {
			jobs <- uint16(port)
		}
	}
	close(jobs)

	wg.Wait()
	close(results)

	<-done

	slices.SortFunc(scannedPorts, func(a, b ScannedPort) int {
		return cmp.Compare(b.Port, a.Port)
	})
	return scannedPorts, nil
}

func (s *Scanner) worker(jobs <-chan uint16, results chan<- ScannedPort, wg *sync.WaitGroup) {
	defer wg.Done()
	dialer := net.Dialer{Timeout: s.config.dialTimeout}

	for port := range jobs {
		address := net.JoinHostPort(s.address.String(), fmt.Sprint(port))
		conn, err := dialer.Dial("tcp", address)
		if err == nil {
			conn.Close()
		}
		result := ScannedPort{Port: port, Opened: err == nil}

		results <- result
	}
}

func (s *Scanner) increaseStats(opened bool) {
	if opened == true {
		s.Stats.Opened++
	} else {
		s.Stats.Closed++
	}
}
