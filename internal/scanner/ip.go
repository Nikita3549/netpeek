package scanner

import (
	"fmt"
	"net"
)

func resolveIP(host string) (net.IP, error) {
	if IP := net.ParseIP(host); IP != nil {
		IPv4 := IP.To4()

		if IPv4 == nil {
			return nil, fmt.Errorf("provided IP is not IPv4")
		}

		return IPv4, nil
	}

	IPs, err := net.LookupHost(host)
	if err != nil || len(IPs) == 0 {
		return nil, fmt.Errorf("failed to resolve host %q: %v", host, err)
	}

	for _, IP := range IPs {
		IPv4 := net.ParseIP(IP).To4()

		if IPv4 != nil {
			return IPv4, nil
		}
	}

	return nil, fmt.Errorf("host %q has no IPv4 addresses (IPv6 not supported)", host)
}
