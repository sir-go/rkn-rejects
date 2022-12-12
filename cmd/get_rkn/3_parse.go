package main

// Unpack a zip dump and parses contained XML files

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/tamerh/xml-stream-parser"
	"golang.org/x/text/encoding/charmap"
)

var (
	reHost *regexp.Regexp
)

// domainIsOk checks if a domain itself is in the white list or of any of it's children is in the whitelist
func domainIsOk(d string) bool {
	if b, s := isItUpDomainOf(d, CFG.WListDomains); b {
		if s == d {
			log.Debugf("`%s` ∈[white-list] -> skip", forcePDecode(d))
		} else {
			log.Debugf(
				"`%s` is upper domain for `%s` ∈[white-list] -> skip",
				forcePDecode(d), forcePDecode(s))
		}
		return false
	}
	return true
}

// processRegElement parses an XML element of dump and toss it to the certain db
func processRegElement(
	elChan chan *xmlparser.XMLElement,
	rdb *redis.Client, wg *sync.WaitGroup) {

	defer wg.Done()
	ctx := context.Background()

	// create new buffer of requests for IPs, domains and check URLs
	buff := NewRBuff("ip_", "domains_", "check_")

	for el := range elChan {
		hsh := el.Attrs["hash"]
		switch el.Attrs["blockType"] {

		// el has blockType=ip attr
		case "ip":

			// parse all of the ipSubnet child elements
			for _, ch := range el.Childs["ipSubnet"] {
				for _, h := range CIDRHosts(ch.InnerText) {
					// make the hash and toss all of IPs of the parsed subnet to the checks
					buff.add("check_", fmt.Sprintf("%s|%s", hsh, h))
					// toss all IPs of the parsed subnet to the IP list
					buff.add("ip_", h)
				}
			}

			// parse all of the ip child elements
			for _, ch := range el.Childs["ip"] {
				for _, h := range CIDRHosts(ch.InnerText) {
					// make the hash and toss all of IPs of the parsed subnet to the checks
					buff.add("check_", fmt.Sprintf("%s|%s", hsh, h))
					// toss all IPs of the parsed subnet to the IP list
					buff.add("ip_", h)
				}
			}

		// el has blockType=domain or blockType=domain-mask attr
		case "domain":
			fallthrough
		case "domain-mask":
			// parse all of the domain child elements
			for _, ch := range el.Childs["domain"] {
				// clean the domain names
				domain := sanitizeDomain(ch.InnerText)
				// check if the domain is allowed
				if domainIsOk(domain) {
					// make the hash and toss the domain to the checks
					buff.add("check_", fmt.Sprintf("%s|%s", hsh, domain))
					//toss the domain to the domains list
					buff.add("domains_", domain)
				}
			}
		// el has blockType=url or anything else attr
		default:
			// parse all of the url child elements
			for _, ch := range el.Childs["url"] {
				// get the hostname from the element content
				subStrings := reHost.FindStringSubmatch(ch.InnerText)
				if len(subStrings) < 2 {
					log.Debugf("[%s] can't get host from url `%s`",
						hsh, subStrings)
					continue
				}

				// got the hostname
				host := subStrings[1]

				// if it an IP address - parse it as a subnet
				// and toss all of parsed addresses to the check and ip lists
				if ip := net.ParseIP(host); ip != nil {
					for _, h := range CIDRHosts(host) {
						buff.add("check_", fmt.Sprintf("%s|%s",
							hsh, ch.InnerText))
						buff.add("ip_", h)
					}
					continue
				}

				// if it's not an IP address then clean the hostname and toss it to the checks and domains lists
				if domain := sanitizeDomain(host); domainIsOk(domain) {
					buff.add("check_", fmt.Sprintf("%s|%s",
						hsh, ch.InnerText))
					buff.add("domains_", domain)
				}
			}
		}

		// if request buffer is full - send it to the redis
		if buff.count > CFG.Parse.Redis.ChunkSize {
			buff.send(rdb, ctx)
		}
	}

	// send the rest of requests to the redis
	if buff.count > 0 {
		buff.send(rdb, ctx)
	}
}

// Parse reads data from an unpacked dump with chunks sized by `size`, decodes it as an XML,
//and parses all of the elements with the progress updating
func Parse(r io.ReadCloser, size uint64) {
	log.Info("start parsing")

	// regexp for get a hostname from the url
	reHost = regexp.MustCompile(`^[a-zA-Z\d]+://([^/\n:\\]+)`)

	// create a new redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			CFG.Parse.Redis.Host,
			CFG.Parse.Redis.Port),
		Password:     CFG.Parse.Redis.Password,
		DB:           CFG.Parse.Redis.Db,
		MaxRetries:   99,
		DialTimeout:  CFG.Parse.Redis.TimeoutConn,
		ReadTimeout:  CFG.Parse.Redis.TimeoutRead,
		WriteTimeout: CFG.Parse.Redis.TimeoutRead,
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Errorln("can't close redis connections", err)
		}
	}()

	// dumps saved in the Windows-1251 encoding
	dec := charmap.Windows1251.NewDecoder()

	// prepare the XML parser
	parser := xmlparser.NewXMLParser(
		bufio.NewReader(dec.Reader(r)),
		"content")
	elChan := parser.Stream()

	// prepare a waiting group for goroutines
	var wg sync.WaitGroup
	waitCh := make(chan struct{})

	// set capacity of the waiting group by the workers amount
	wg.Add(CFG.Parse.Redis.Workers)

	// create a ticker for the progressbar updating
	tick := time.NewTicker(CFG.Parse.ProgressPollTimeout)

	// run by-element parsing workers
	go func() {
		for i := CFG.Parse.Redis.Workers; i > 0; i-- {
			go processRegElement(elChan, rdb, &wg)
		}
		wg.Wait()
		close(waitCh)
	}()

	var (
		progress uint64
		err      error
	)

	// a loop for goroutines processing
	for {
		select {

		// all workers are done
		case <-waitCh:
			if progress < 100 {
				log.Info("100%")
			}

			ctx := context.Background()

			// rename the temporary list on the redis (remove a _ prefix)
			rKeys := rdb.Keys(ctx, "*_")
			if err = rKeys.Err(); err != nil {
				log.Panicln("redis get keys", err)
			}
			for _, k := range rKeys.Val() {
				err = rdb.Rename(ctx, k, strings.TrimSuffix(k, "_")).Err()
				if err != nil {
					log.Panicln("rename redis keys", err)
				}
			}
			return

		// update the progress by the ticker
		case <-tick.C:
			progress = parser.TotalReadSize / (size / 100.0)
			if progress > 100 {
				progress = 100
			}
			log.Infof("%d%%", progress)
		}
	}
}
