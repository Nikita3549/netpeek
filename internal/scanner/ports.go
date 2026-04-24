package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

type PortRange struct {
	Start uint16
	End   uint16
}

type totalPorts int

const minPort = 1
const maxPort = 65535

// Allowed format: all or 80-1024,80-81,443
func parsePorts(ports string) ([]PortRange, totalPorts, error) {
	var total totalPorts

	if ports == "all" {
		return []PortRange{
			{
				Start: minPort,
				End:   maxPort,
			}}, maxPort, nil
	}

	portsList := strings.Split(ports, ",")
	portRanges := make([]PortRange, 0, len(portsList))

	for _, portRange := range portsList {
		isSingle := !strings.Contains(portRange, "-")

		if isSingle {
			portUint16, err := toUint16Port(portRange)

			if err != nil {
				return nil, 0, err
			}

			total++
			portRanges = append(portRanges, PortRange{Start: portUint16, End: portUint16})
		} else {
			portsString := strings.Split(portRange, "-")

			if len(portsString) != 2 {
				return nil, 0, fmt.Errorf("invalid ports range: %q", portsString)
			}

			start, err := toUint16Port(portsString[0])
			if err != nil {
				return nil, 0, err
			}

			end, err := toUint16Port(portsString[1])
			if err != nil {
				return nil, 0, err
			}

			if start > end {
				return nil, 0, fmt.Errorf("start port can't be greater than end port: %q", portRange)
			}

			total += totalPorts(end - start + 1)
			portRanges = append(portRanges, PortRange{Start: start, End: end})
		}
	}

	return portRanges, total, nil
}

func toUint16Port(port string) (uint16, error) {
	portInt, err := strconv.Atoi(strings.TrimSpace(port))

	if err != nil {
		return 0, fmt.Errorf("invalid ports syntax: %q", port)
	}

	if portInt < minPort || portInt > maxPort {
		return 0, fmt.Errorf("port %v out of %v-%v range", portInt, minPort, maxPort)
	}

	return uint16(portInt), nil
}
