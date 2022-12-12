package main

import (
	"net"
	"strings"

	"github.com/c-robinson/iplib"
	log "github.com/sirupsen/logrus"
)

// ipInSubnets checks if the IP address is in the subnets slice and returns a boolean (found or not)
//and the subnet as a string where this IP address has been found
func ipInSubnets(ip net.IP, subnets []Net) (bool, string) {
	for _, sn := range subnets {
		if sn.Contains(ip) {
			return true, sn.String()
		}
	}
	return false, ""
}

// ipInBogusSubnet check if the defined in the configuration bogus networks contain the given IP address
func ipInBogusSubnet(ip net.IP) bool {
	if b, sn := ipInSubnets(ip, CFG.Parse.BogusIp.Subnets); b {
		log.Warnf("%s in bogus subnet %s -> skip", ip, sn)
		return b
	}
	return false
}

// CIDRHosts composites the given CIDR to a slice of hosts (IP strings)
// "10.0.0.12/30 -> 10.0.0.12/32, 10.0.0.13/32, 10.0.0.14/32, 10.0.0.15/32,
func CIDRHosts(cidr string) (hosts []string) {
	// it's IP address
	if !strings.ContainsRune(cidr, '/') {
		if ip := net.ParseIP(cidr); ip != nil && !ipInBogusSubnet(ip) {
			return []string{cidr}
		}
		return
	}

	// it's Subnet address
	lIp, lNet, err := iplib.ParseCIDR(cidr)
	if err != nil || lIp == nil || lNet == nil {
		log.Warnln(cidr, "can't parse CIDR -> skip")
		return
	}

	prefixLen, _ := lNet.Mask().Size()
	if prefixLen < CFG.Parse.BogusIp.MinMask {
		log.Warnln(cidr, "to wide mask -> skip")
		return
	}

	if prefixLen == 32 {
		if ipInBogusSubnet(lIp) {
			return
		}
		return []string{lIp.String()}
	}

	for _, hostIP := range lNet.(iplib.Net4).Enumerate(0, 0) {
		if ipInBogusSubnet(hostIP) {
			continue
		}
		hosts = append(hosts, hostIP.String())
	}

	return
}
