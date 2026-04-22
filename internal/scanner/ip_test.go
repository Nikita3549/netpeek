package scanner

import (
	"net/netip"
	"testing"
)

func TestResolveIp(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		wantErr     bool
		wantAddress netip.Addr
	}{
		{
			name:        "valid IPv4",
			host:        "8.8.8.8",
			wantErr:     false,
			wantAddress: generateAddress("8.8.8.8"),
		},
		{
			name:        "valid IPv6",
			host:        "::1",
			wantErr:     false,
			wantAddress: generateAddress("::1"),
		},
		{
			name:        "valid hostname",
			host:        "localhost",
			wantErr:     false,
			wantAddress: generateAddress("127.0.0.1"),
		},
		{
			name:        "invalid host",
			host:        "test.xyz",
			wantErr:     true,
			wantAddress: netip.Addr{},
		},
		{
			name:        "empty string",
			host:        "",
			wantErr:     true,
			wantAddress: netip.Addr{},
		},
		{
			name:        "garbage input",
			host:        "not an ip!!!",
			wantErr:     true,
			wantAddress: netip.Addr{},
		},
	}

	for _, tt := range tests {
		got, err := resolveIP(tt.host)

		t.Run(tt.name, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveIP(%q) error = %v, wantErr = %v", tt.host, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !got.IsValid() {
					t.Errorf("resolveIp(%q) got = nil, want = %v", tt.host, tt.wantAddress)
					return
				}

				if got != tt.wantAddress {
					t.Errorf("resolveIp(%q) got = %v, want = %v", tt.host, got, tt.wantAddress)
					return
				}
			}
		})
	}
}

func generateAddress(address string) netip.Addr {
	addr, _ := netip.ParseAddr(address)

	return addr
}
