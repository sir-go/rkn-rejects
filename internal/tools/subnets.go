package tools

import (
	"net"
)

// IpInSubnets check if the IP address is found in the subnets
// return the boolean (found IP or not) and the subnet that contains the given IP as a string
func IpInSubnets(ip net.IP, subnets []*net.IPNet) (bool, string) {
	for _, sn := range subnets {
		if sn.Contains(ip) {
			return true, sn.String()
		}
	}
	return false, ""
}
