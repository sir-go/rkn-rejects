package main

// Some tools for domain names working with

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/idna"
)

// isASCII checks if the whole string is ASCII encoded
// example -> true
// пример -> false
// some-пример -> false
func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

// forcePUnicode encodes a domain as punycode if it has certain signature
// сайт.рф -> xn--80aswg.xn--p1ai
// example.com -> example.com
func forcePUnicode(d string) (s string) {
	var err error
	if !isASCII(d) {
		if s, err = idna.Punycode.ToASCII(d); err != nil {
			log.Warn(err)
			return d
		}
	} else {
		s = d
	}
	return
}

// forcePDecode decodes a domain as punycode if it has certain signature
// xn--80aswg.xn--p1ai -> сайт.рф
// example.com -> example.com
func forcePDecode(d string) (s string) {
	if !strings.HasPrefix(d, "xn--") {
		return d
	}
	var err error
	if s, err = idna.Punycode.ToUnicode(d); err != nil {
		log.Warn(err)
		return d
	}
	return
}

// sanitizeDomain clears a domain name
//www.domain.com -> domain.com
//*.sub.domain.net -> sub.domain.net
//domain.com/?params=444 -> domain.com
func sanitizeDomain(d string) (s string) {
	s = forcePUnicode(d)
	s = strings.TrimPrefix(s, "www.")
	s = strings.TrimPrefix(s, "*.")
	if dzI := strings.IndexAny(s, "#/?:+"); dzI != -1 {
		s = s[:dzI]
	}
	s = strings.TrimRightFunc(s, func(r rune) bool { return r == '.' })
	return
}

// isItUpDomainOf checks if any of the given subdomains is a child of the given domain.
// Returns a boolean (is the domain is the parent) and subdomain string
func isItUpDomainOf(domain string, subDomains []string) (bool, string) {
	for _, sd := range subDomains {
		if strings.HasSuffix(sd, "."+domain) || sd == domain {
			return true, sd
		}
	}
	return false, ""
}
