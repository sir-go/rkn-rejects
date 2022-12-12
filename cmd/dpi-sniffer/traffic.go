package main

import (
	"regexp"
	"strings"

	"github.com/florianl/go-nfqueue"
	"github.com/go-redis/redis/v8"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	log "github.com/sirupsen/logrus"

	"rkn-rejects/internal/tools"
)

//goland:noinspection SpellCheckingInspection
var (
	// ReLooksLikeDomain matches strings tooks like domain name
	ReLooksLikeDomain = regexp.MustCompile(`([a-zA-Z]+[a-zA-Z0-9._-]*\.)+([a-zA-Z]+[a-zA-Z0-9._-]*)`)

	// ReHostname matches strings tooks like host name
	ReHostname = regexp.MustCompile(`(?i)host:\s([^:/?#\s]+).*\s`)

	// ReHTTP matches strings tooks like HTTP header
	ReHTTP = regexp.MustCompile(`^(GET|POST|PUT|PATCH|DELETE|TRACE|CONNECT|HEAD|OPTIONS)\s`)
)

// snatchRegexp gets a host or domain name from the given slice
// pos argument sets a position in the matched results
func snatchRegexp(b []byte, re *regexp.Regexp, pos int) (s string) {
	if reResult := re.FindSubmatch(b); reResult != nil && len(reResult) > 1 {
		return strings.TrimRightFunc(string(reResult[pos]), func(r rune) bool { return r == '.' })
	}
	return ""
}

// packetsHook - nf queue packet processor
func packetsHook(a nfqueue.Attribute, rdb *redis.Client) bool {

	// nf verdicts mapping
	const (
		accept = false
		reject = true
	)

	// if -x (rejects all) flag is set - reject the packet without any checks
	if CFG.RejectAll {
		log.Warn("R-ALL: any =>X any")
		return reject
	}

	// parse the packet as IPv4
	p := gopacket.NewPacket(*a.Payload, layers.LayerTypeIPv4, gopacket.DecodeOptions{
		Lazy:                     true,
		NoCopy:                   true,
		SkipDecodeRecovery:       true,
		DecodeStreamsAsDatagrams: true,
	})

	// check dst IPv4
	srcIP, dstIP := p.NetworkLayer().NetworkFlow().Endpoints()
	log.Debugf("IP  : %-15s <-> %s", srcIP, dstIP)

	if p.TransportLayer() == nil {
		log.Debugln("transport is nil -> X")
		return reject
	}

	if p.TransportLayer().LayerPayload() == nil {
		log.Debugln("transport payload is nil -> X")
		return reject
	}

	// check TLS
	//log.Debugln("check TLS")
	var hostname string
	if hostname, _ = GetSNIForced(p.TransportLayer().LayerPayload()); len(hostname) > 1 {
		if !tools.IsHostDenied(hostname, rdb, CFG.Redis.SetKey) {
			log.Debugf("TLS : %-15s ==> %s", srcIP, hostname)
			return accept
		}
		log.Infof("TLS : %-15s =>X %s", srcIP, hostname)
		return reject
	}

	// check HTTP hostname
	//log.Debugln("check hostname")
	if hostname = snatchRegexp(*a.Payload, ReHostname, 1); hostname != "" {
		if !tools.IsHostDenied(hostname, rdb, CFG.Redis.SetKey) {
			log.Debugf("HTTP: %-15s ==> %s", srcIP, hostname)
			return accept
		}
		log.Infof("HTTP: %-15s =>X %s", srcIP, hostname)
		return reject
	}

	// has payload a http method?
	//log.Debugln("check http method")
	if ReHTTP.Match(*a.Payload) {
		log.Infof("noHN: %-15s =>X %s", srcIP, dstIP)
		return reject
	}

	// DPI seek anything looks like domain name
	//log.Debugln("check DPI")
	if hostname = snatchRegexp(*a.Payload, ReLooksLikeDomain, 0); hostname != "" {
		if !tools.IsHostDenied(hostname, rdb, CFG.Redis.SetKey) {
			log.Debugf("DPI : %-15s ==> %s", srcIP, hostname)
			return accept
		}
		log.Infof("DPI : %-15s =>X %s", srcIP, hostname)
		return reject
	}

	log.Debugf("??? : %-15s ==> %s", srcIP, dstIP)
	return accept
}
