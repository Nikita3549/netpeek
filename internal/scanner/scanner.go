package scanner

import (
	"cmp"
	"fmt"
	"net"
	"slices"
	"sync"
	"time"
)

const dialTimeout = 500 * time.Millisecond
const workerCount = 1024

type Scanner struct {
	ip    net.IP
	ports []PortRange
}

type ScannedPort struct {
	Port   uint16
	Opened bool
}

func NewScanner(host, ports string) (*Scanner, error) {
	ip, err := resolveIP(host)
	if err != nil {
		return nil, err
	}

	portRanges, err := parsePorts(ports)
	if err != nil {
		return nil, err
	}

	return &Scanner{
		ports: portRanges,
		ip:    ip,
	}, nil
}

func (s *Scanner) Scan() ([]ScannedPort, error) {
	jobs := make(chan uint16, workerCount)
	results := make(chan ScannedPort)
	scannedPorts := make([]ScannedPort, 0)
	var wg sync.WaitGroup

	// Worker
	for i := 0; i <= workerCount; i++ {
		wg.Add(1)
		go s.worker(jobs, results, &wg)
	}

	// Collector
	done := make(chan bool)
	go func() {
		for result := range results {
			scannedPorts = append(scannedPorts, result)
		}
		done <- true
	}()

	// Producer
	for _, portRange := range s.ports {
		for port := portRange.Start; port <= portRange.End; port++ {
			jobs <- port
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
	dialer := net.Dialer{Timeout: dialTimeout}

	for port := range jobs {
		address := net.JoinHostPort(s.ip.String(), fmt.Sprint(port))
		conn, err := dialer.Dial("tcp", address)
		if err == nil {
			conn.Close()
		}
		result := ScannedPort{Port: port, Opened: err == nil}

		results <- result
	}
}
