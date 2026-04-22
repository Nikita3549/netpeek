package scanner

import (
	"fmt"
	"net"
	"net/netip"
)

func resolveIP(host string) (netip.Addr, error) {
	if address, err := netip.ParseAddr(host); err == nil {
		return address, nil
	}

	IPs, err := net.LookupHost(host)
	if err != nil || len(IPs) == 0 {
		return netip.Addr{}, fmt.Errorf("failed to resolve host %q: %v", host, err)
	}

	var fallback netip.Addr
	for _, IP := range IPs {
		address, _ := netip.ParseAddr(IP)

		if address.Is4() {
			return address, nil
		}
		if address.Is6() {
			fallback = address
		}
	}
	if fallback.IsValid() {
		return fallback, nil
	}

	return netip.Addr{}, fmt.Errorf("couldn't resolve host %q", host)
}
