package cli

import (
	"net"
	"strconv"
	"testing"
)

func TestScanCommand_RequiresHostAndPortFlags(t *testing.T) {
	rootCmd.SetArgs([]string{"scan"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for missing required flags")
	}
}

func TestScanCommand_SucceedsForReachableLocalPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	defer ln.Close()

	_, portString, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort() failed: %v", err)
	}

	rootCmd.SetArgs([]string{
		"scan",
		"--host", "127.0.0.1",
		"--port", portString,
		"--concurrency", "1",
		"--timeout", "100",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected successful scan, got error: %v", err)
	}
}

func TestScanCommand_ReturnsErrorForInvalidPort(t *testing.T) {
	rootCmd.SetArgs([]string{
		"scan",
		"--host", "127.0.0.1",
		"--port", "70000",
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for invalid port")
	}
}

func TestScanCommand_SetsConcurrencyAndTimeoutFromFlags(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	defer ln.Close()

	_, portString, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort() failed: %v", err)
	}

	wantConcurrency := 3
	wantTimeout := 150

	rootCmd.SetArgs([]string{
		"scan",
		"-H", "127.0.0.1",
		"-p", portString,
		"-c", strconv.Itoa(wantConcurrency),
		"-t", strconv.Itoa(wantTimeout),
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected successful scan, got error: %v", err)
	}

	if concurrency != wantConcurrency {
		t.Fatalf("concurrency = %d, want %d", concurrency, wantConcurrency)
	}

	if timeout != wantTimeout {
		t.Fatalf("timeout = %d, want %d", timeout, wantTimeout)
	}
}
