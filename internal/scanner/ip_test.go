package scanner

import (
	"net"
	"testing"
)

func TestResolveIp(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
		wantIP  net.IP
	}{
		{
			name:    "valid IPv4",
			host:    "8.8.8.8",
			wantErr: false,
			wantIP:  net.ParseIP("8.8.8.8"),
		},
		{
			name:    "valid IPv6",
			host:    "::1",
			wantErr: true,
			wantIP:  net.ParseIP("::1"),
		},
		{
			name:    "valid hostname",
			host:    "localhost",
			wantErr: false,
			wantIP:  net.ParseIP("127.0.0.1"),
		},
		{
			name:    "invalid host",
			host:    "test.xyz",
			wantErr: true,
			wantIP:  nil,
		},
		{
			name:    "empty string",
			host:    "",
			wantErr: true,
			wantIP:  nil,
		},
		{
			name:    "garbage input",
			host:    "not an ip!!!",
			wantErr: true,
			wantIP:  nil,
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
				if got == nil {
					t.Errorf("resolveIp(%q) got = nil, want = %v", tt.host, tt.wantIP)
					return
				}

				if !got.Equal(tt.wantIP) {
					t.Errorf("resolveIp(%q) got = %v, want = %v", tt.host, got, tt.wantIP)
					return
				}
			}
		})
	}
}
