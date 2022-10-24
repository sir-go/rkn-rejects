package main

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

func processRegElement(
	elChan chan *xmlparser.XMLElement,
	rdb *redis.Client, wg *sync.WaitGroup) {

	defer wg.Done()
	ctx := context.Background()
	buff := NewRBuff("ip_", "domains_", "check_")

	for el := range elChan {
		hsh := el.Attrs["hash"]
		switch el.Attrs["blockType"] {

		case "ip":
			for _, ch := range el.Childs["ipSubnet"] {
				for _, h := range CIDRHosts(ch.InnerText) {
					buff.add("check_", fmt.Sprintf("%s|%s", hsh, h))
					buff.add("ip_", h)
				}
			}
			for _, ch := range el.Childs["ip"] {
				for _, h := range CIDRHosts(ch.InnerText) {
					buff.add("check_", fmt.Sprintf("%s|%s", hsh, h))
					buff.add("ip_", h)
				}
			}
		case "domain":
			fallthrough
		case "domain-mask":
			for _, ch := range el.Childs["domain"] {
				domain := sanitizeDomain(ch.InnerText)
				if domainIsOk(domain) {
					buff.add("check_", fmt.Sprintf("%s|%s", hsh, domain))
					buff.add("domains_", domain)
				}
			}
		default:
			for _, ch := range el.Childs["url"] {
				subStrings := reHost.FindStringSubmatch(ch.InnerText)
				if len(subStrings) < 2 {
					log.Debugf("[%s] can't get host from url `%s`",
						hsh, subStrings)
					continue
				}

				host := subStrings[1]
				if ip := net.ParseIP(host); ip != nil {
					for _, h := range CIDRHosts(host) {
						buff.add("check_", fmt.Sprintf("%s|%s",
							hsh, ch.InnerText))
						buff.add("ip_", h)
					}
					continue
				}

				if domain := sanitizeDomain(host); domainIsOk(domain) {
					buff.add("check_", fmt.Sprintf("%s|%s",
						hsh, ch.InnerText))
					buff.add("domains_", domain)
				}
			}
		}

		if buff.count > CFG.Parse.Redis.ChunkSize {
			buff.send(rdb, ctx)
		}
	}
	if buff.count > 0 {
		buff.send(rdb, ctx)
	}
}

func Parse(r io.ReadCloser, size uint64) {
	log.Info("start parsing")
	reHost = regexp.MustCompile(`^[a-zA-Z\d]+://([^/\n:\\]+)`)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			CFG.Parse.Redis.Host,
			CFG.Parse.Redis.Port),
		Password:     CFG.Parse.Redis.Password,
		DB:           CFG.Parse.Redis.Db,
		MaxRetries:   99,
		DialTimeout:  CFG.Parse.Redis.TimeoutConn.Duration,
		ReadTimeout:  CFG.Parse.Redis.TimeoutRead.Duration,
		WriteTimeout: CFG.Parse.Redis.TimeoutRead.Duration,
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Errorln("can't close redis connections", err)
		}
	}()

	dec := charmap.Windows1251.NewDecoder()

	parser := xmlparser.NewXMLParser(
		bufio.NewReader(dec.Reader(r)),
		"content")
	elChan := parser.Stream()

	var wg sync.WaitGroup
	waitCh := make(chan struct{})

	wg.Add(CFG.Parse.Redis.Workers)
	tick := time.NewTicker(CFG.Parse.ProgressPollTimeout.Duration)

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
	for {
		select {
		case <-waitCh:
			if progress < 100 {
				log.Info("100%")
			}

			ctx := context.Background()
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
		case <-tick.C:
			progress = parser.TotalReadSize / (size / 100.0)
			if progress > 100 {
				progress = 100
			}
			log.Infof("%d%%", progress)
		}
	}
}
