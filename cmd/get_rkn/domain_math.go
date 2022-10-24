package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/idna"
)

func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

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

func isItUpDomainOf(domain string, subDomains []string) (bool, string) {
	for _, sd := range subDomains {
		if strings.HasSuffix(sd, "."+domain) || sd == domain {
			return true, sd
		}
	}
	return false, ""
}
