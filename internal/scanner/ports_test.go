package scanner

import (
	"slices"
	"testing"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		name           string
		ports          string
		wantErr        bool
		wantPortRanges []PortRange
	}{
		{
			name:    "all ports",
			ports:   "all",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: minPort, End: maxPort},
			},
		},
		{
			name:    "single port",
			ports:   "80",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: 80, End: 80},
			},
		},
		{
			name:    "port range",
			ports:   "80-443",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: 80, End: 443},
			},
		},
		{
			name:    "multiple ports comma",
			ports:   "22,80,443",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: 22, End: 22},
				{Start: 80, End: 80},
				{Start: 443, End: 443},
			},
		},
		{
			name:    "mixed ports and ranges",
			ports:   "22,80-100,443",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: 22, End: 22},
				{Start: 80, End: 100},
				{Start: 443, End: 443},
			},
		},
		{
			name:    "invalid string",
			ports:   "abc",
			wantErr: true,
		},
		{
			name:    "reversed range",
			ports:   "443-80",
			wantErr: true,
		},
		{
			name:    "out of range port",
			ports:   "99999",
			wantErr: true,
		},
		{
			name:    "empty string",
			ports:   "",
			wantErr: true,
		},
		{
			name:    "port zero",
			ports:   "0",
			wantErr: true,
		},
		{
			name:    "max port",
			ports:   "65535",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: 65535, End: 65535},
			},
		},
		{
			name:    "port above max",
			ports:   "65536",
			wantErr: true,
		},
		{
			name:    "range touching boundaries",
			ports:   "1-65535",
			wantErr: false,
			wantPortRanges: []PortRange{
				{Start: 1, End: 65535},
			},
		},
		{
			name:    "trailing comma",
			ports:   "80,",
			wantErr: true,
		},
		{
			name:    "double dash",
			ports:   "80--443",
			wantErr: true,
		},
		{
			name:    "spaces",
			ports:   "80, 443",
			wantErr: false,
			wantPortRanges: []PortRange{
				{
					Start: 80,
					End:   80,
				},
				{
					Start: 443,
					End:   443,
				},
			},
		},
		{
			name:    "wrong first port in range",
			ports:   "0-443",
			wantErr: true,
		},
		{
			name:    "wrong second port in range",
			ports:   "1-65536",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePorts(tt.ports)

			if (err != nil) != tt.wantErr {
				t.Errorf("parsePorts(%q) err = %v wantErr = %v", tt.ports, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("parsePorts(%q) got = nil, wantPortRanges = %v", tt.ports, tt.wantPortRanges)
					return
				}

				if !slices.Equal(got, tt.wantPortRanges) {
					t.Errorf("parsePorts(%q) got = %v, wantPortRanges = %v", tt.ports, got, tt.wantPortRanges)
				}
			}
		})
	}
}
