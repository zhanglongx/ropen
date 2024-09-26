package pkg

import (
	"fmt"
	"net"
	"strings"
)

const (
	SUBNET_MASK = "255.255.255.0"
)

func getIPs(prefer []string) ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []string

	for _, iface := range interfaces {
		debug("checking interface %v", iface.Name)

		if !isPysicalInterface(iface) {
			debug("skipping interface %v", iface.Name)
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			debug("failed to get addresses for interface %v: %v", iface.Name, err)
			continue
		}

		if len(addrs) == 0 {
			continue
		}

		for _, addr := range addrs {
			debug("found address %v for interface %v", addr, iface.Name)

			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				debug("failed to parse address: %v", err)
				continue
			}

			if len(prefer) == 0 {
				ips = append(ips, ip.String())
				continue
			} else {
				for _, p := range prefer {
					same, err := isSameSubnet(ip.String(), p)
					if err != nil {
						debug("failed to compare addresses: %v", err)
						continue
					}

					if same {
						ips = append(ips, ip.String())
						break
					}
				}
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no valid IP addresses found")
	}

	return ips, nil
}

func isPysicalInterface(iface net.Interface) bool {
	return !(strings.HasPrefix(iface.Name, "docker") ||
		strings.HasPrefix(iface.Name, "veth")) &&
		iface.Flags&net.FlagLoopback == 0 &&
		iface.Flags&net.FlagPointToPoint == 0 &&
		iface.Flags&net.FlagUp != 0
}

func isSameSubnet(ip1, ip2 string) (bool, error) {
	parsedIP1 := net.ParseIP(ip1)
	parsedIP2 := net.ParseIP(ip2)
	if parsedIP1 == nil || parsedIP2 == nil {
		return false, fmt.Errorf("invalid IP address")
	}

	parsedMask := net.IPMask(net.ParseIP(SUBNET_MASK).To4())
	if parsedMask == nil {
		return false, fmt.Errorf("invalid subnet mask")
	}

	network1 := parsedIP1.Mask(parsedMask)
	network2 := parsedIP2.Mask(parsedMask)

	return network1.Equal(network2), nil
}
