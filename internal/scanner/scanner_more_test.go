package scanner

import (
	"fmt"
	"net"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestNewScanner_UsesDefaultConfig(t *testing.T) {
	s, err := NewScanner("127.0.0.1", "80")
	if err != nil {
		t.Fatalf("NewScanner() failed: %v", err)
	}

	if s.config.workerCount != DefaultWorkerCount {
		t.Fatalf("workerCount = %d, want %d", s.config.workerCount, DefaultWorkerCount)
	}

	if s.config.dialTimeout != DefaultDialTimeout {
		t.Fatalf("dialTimeout = %v, want %v", s.config.dialTimeout, DefaultDialTimeout)
	}
}

func TestNewScanner_AppliesCustomConfig(t *testing.T) {
	wantWorkers := 4
	wantTimeout := 25 * time.Millisecond

	s, err := NewScanner(
		"127.0.0.1",
		"80",
		WithWorkerCount(wantWorkers),
		WithDialTimeout(wantTimeout),
	)
	if err != nil {
		t.Fatalf("NewScanner() failed: %v", err)
	}

	if s.config.workerCount != wantWorkers {
		t.Fatalf("workerCount = %d, want %d", s.config.workerCount, wantWorkers)
	}

	if s.config.dialTimeout != wantTimeout {
		t.Fatalf("dialTimeout = %v, want %v", s.config.dialTimeout, wantTimeout)
	}
}

func TestScanner_IncreaseStats(t *testing.T) {
	s := &Scanner{}

	s.increaseStats(true)
	s.increaseStats(false)
	s.increaseStats(false)

	if s.Stats.Opened != 1 {
		t.Fatalf("opened = %d, want 1", s.Stats.Opened)
	}

	if s.Stats.Closed != 2 {
		t.Fatalf("closed = %d, want 2", s.Stats.Closed)
	}
}

func TestScanner_ScanSortsResultsDescendingByPort(t *testing.T) {
	ln1, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start first listener: %v", err)
	}
	defer ln1.Close()

	ln2, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start second listener: %v", err)
	}
	defer ln2.Close()

	_, port1String, err := net.SplitHostPort(ln1.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort() for first listener failed: %v", err)
	}

	_, port2String, err := net.SplitHostPort(ln2.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort() for second listener failed: %v", err)
	}

	ports := fmt.Sprintf("%s,%s", port1String, port2String)
	s, err := NewScanner("127.0.0.1", ports, WithWorkerCount(2))
	if err != nil {
		t.Fatalf("NewScanner() failed: %v", err)
	}

	p := NewPrinter(s.Stats.Total)
	results, err := s.Scan(p)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}

	gotPorts := []uint16{results[0].Port, results[1].Port}
	wantPorts := make([]uint16, 0, 2)

	port1, err := strconv.Atoi(port1String)
	if err != nil {
		t.Fatalf("Atoi() for first port failed: %v", err)
	}
	port2, err := strconv.Atoi(port2String)
	if err != nil {
		t.Fatalf("Atoi() for second port failed: %v", err)
	}

	wantPorts = append(wantPorts, uint16(port1), uint16(port2))
	slices.SortFunc(wantPorts, func(a, b uint16) int {
		if a > b {
			return -1
		}
		if a < b {
			return 1
		}

		return 0
	})

	if !slices.Equal(gotPorts, wantPorts) {
		t.Fatalf("ports order = %v, want %v", gotPorts, wantPorts)
	}

	if !results[0].Opened || !results[1].Opened {
		t.Fatalf("expected both ports to be open, got %+v", results)
	}
}
