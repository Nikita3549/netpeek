package scanner

import (
	"net"
	"testing"
)

func TestScanner_Integration(t *testing.T) {
	p := NewPrinter(100)

	t.Run("Scan Opened Port", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Errorf("could not start listener: %v", err)
		}

		address := ln.Addr().String()
		host, port, _ := net.SplitHostPort(address)

		go func() {
			conn, err := ln.Accept()

			if err == nil {
				conn.Close()
			}
		}()
		defer ln.Close()

		scanner, err := NewScanner(host, port)
		if err != nil {
			t.Fatalf("NewScanner(%v, %v) failed: %v", host, port, err)
		}

		results, err := scanner.Scan(p)
		if err != nil {
			t.Fatalf("Scan() failed: %v", err)
			return
		}

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].Opened {
			t.Errorf("port %s: expected status Opened=true, but got false", port)
		}
	})
	t.Run("Scan Closed Port", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Errorf("could not start listener: %v", err)
		}

		address := ln.Addr().String()
		host, port, _ := net.SplitHostPort(address)

		ln.Close()

		scanner, err := NewScanner(host, port)
		if err != nil {
			t.Fatalf("NewScanner(%v, %v) failed: %v", host, port, err)
		}

		results, err := scanner.Scan(p)
		if err != nil {
			t.Fatalf("Scan() failed: %v", err)
			return
		}

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Opened {
			t.Errorf("port %s: expected status Opened=false, but got true", port)
		}
	})
}
