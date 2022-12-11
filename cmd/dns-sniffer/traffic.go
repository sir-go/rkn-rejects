package main

import (
	"context"
	"net"
	"time"

	"github.com/florianl/go-nfqueue"
	"github.com/go-redis/redis/v8"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	log "github.com/sirupsen/logrus"

	"rkn-rejects/internal/fw"
)

// ipInSubnets check if the IP address is found in the subnets
// return the boolean (found IP or not) and the subnet that contains the given IP as a string
func ipInSubnets(ip net.IP, subnets []*net.IPNet) (bool, string) {
	for _, sn := range subnets {
		if sn.Contains(ip) {
			return true, sn.String()
		}
	}
	return false, ""
}

// getUpperDomains returns a list of all upper domains of the given hostname
func getUpperDomains(d string) (res []string) {
	var dotIndexes []int
	for idx, r := range d {
		if r == '.' {
			dotIndexes = append(dotIndexes, idx)
		}
	}
	for i := len(dotIndexes) - 2; i > -1; i-- {
		res = append(res, d[dotIndexes[i]+1:])
	}
	return append(res, d)
}

// isHostDenied checks the hostname and all of the upper domains if any of them is in the denied list
func isHostDenied(h string, rdb *redis.Client) bool {
	ctx := context.Background()
	for _, ud := range getUpperDomains(h) {
		rRes := rdb.SIsMember(ctx, CFG.Redis.SetKey, ud)
		if err := rRes.Err(); err != nil {
			log.Errorln("redis sismember", CFG.Redis.SetKey, err)
			return false
		}
		if rRes.Val() {
			return true
		}
	}
	return false
}

// queue processing function
// packetsHook gets DNS answer packet from NF, parses it, checks if the hostname denied,
// then checks if resolved IP addresses in the allowed IP addresses cache,
// if the host is denied and IP in the allowed cache - remove IP fromm the cache,
// if the host isn't denied and the cache hasn't it's IP - add IP to the cache
func packetsHook(a nfqueue.Attribute, rdb *redis.Client, ac *AllowsCache) {
	p := gopacket.NewPacket(*a.Payload, layers.LayerTypeIPv4,
		gopacket.DecodeOptions{
			Lazy:                     true,
			NoCopy:                   true,
			SkipDecodeRecovery:       true,
			DecodeStreamsAsDatagrams: true,
		})

	// parse the packet as DNS answer
	l7 := p.ApplicationLayer()
	if l7 == nil {
		return
	}
	if !l7.LayerType().Contains(layers.LayerTypeDNS) {
		return
	}
	dns := l7.(*layers.DNS)

	if dns.Answers == nil || len(dns.Answers) == 0 ||
		dns.Questions == nil || len(dns.Questions) == 0 {
		return
	}

	questedName := string(dns.Questions[0].Name)

	// checks if the hostname is denied
	isDenied := isHostDenied(questedName, rdb)

	// check if IP in bogus networks, add to allowed cache if host isn't denied, and remove if it is
	for _, answer := range dns.Answers {
		if answer.Type != layers.DNSTypeA {
			continue
		}

		// if the IP in bogus networks - skip it
		if bogus, bNet := ipInSubnets(answer.IP, BogusSubnets); bogus {
			log.Debugf("DNS Q: %s A: bogus IP: %s [in %s]",
				questedName, answer.IP.String(), bNet)
			continue
		}

		el := fw.NfSetElement{
			Ip:      answer.IP.To4().String(),
			Timeout: answer.TTL,
			Comment: questedName,
		}

		// check if cache already has the IP
		inCache := ac.Has(el.Ip)
		if isDenied {
			if inCache {

				// host is denied and the IP is in the allowed addresses cache - remove it

				log.Infoln("bad : ", questedName)
				if err := fw.Del(CFG.Nf.Table, CFG.Nf.SetName, el); err != nil {
					log.Error(err)
				}
				if err := ac.Del(el.Ip); err != nil {
					log.Error(err)
				}
			}
		} else {
			if !inCache {

				// host is not denied and the IP is not in the allowed addresses cache - add it

				log.Infoln("good: ", questedName)
				if err := fw.Add(CFG.Nf.Table, CFG.Nf.SetName, el); err != nil {
					log.Error(err)
				}
				err := ac.Add(el.Ip, time.Now().Unix()+int64(el.Timeout))
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
