package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type LogFormat struct{}

func (l *LogFormat) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%-29s [%s] %-15s: %v\n",
		entry.Time.Format("2006/01/02 15:04:05.00.999999"),
		entry.Level.String()[0:1],
		fmt.Sprint(path.Base(entry.Caller.File), ":", entry.Caller.Line),
		entry.Message),
	), nil
}

var (
	CFG          *Cfg
	BogusSubnets []*net.IPNet
	td           []func()
)

func InitInterrupt() {
	log.Warn("-- start --")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func(c chan os.Signal) {
		for range c {
			for _, tdf := range td {
				tdf()
			}
			log.Warn("-- stop --")
			os.Exit(137)
		}
	}(c)
}

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&LogFormat{})
	log.SetLevel(log.DebugLevel)
	InitInterrupt()
	CFG = initConfig()

	if CFG.LogLevel != "debug" {
		logLevel, err := log.ParseLevel(CFG.LogLevel)
		if err != nil {
			log.Panic("parsing LogLevel error")
		}
		log.SetLevel(logLevel)
	}

	for _, ipS := range []string{
		"0.0.0.0/8",      // 0.0.0.0     - 0.255.255.255
		"10.0.0.0/8",     // 10.0.0.0    - 10.255.255.255
		"14.0.0.0/8",     // 14.0.0.0    - 10.255.255.255
		"169.254.0.0/16", // 169.254.0.0 - 169.254.255.255
		"127.0.0.0/8",    // 127.0.0.0   - 127.255.255.255
		"192.168.0.0/16", // 192.168.0.0 - 192.168.255.255
		"172.16.0.0/12",  // 172.16.0.0  - 172.31.255.255
		"192.0.2.0/24",   // 192.0.2.0   - 192.0.2.255
		"224.0.0.0/3",    // 224.0.0.0   - 255.255.255.255
	} {
		_, n, err := net.ParseCIDR(ipS)
		if err != nil {
			log.Panicln("can't parse CIDR: ", ipS, err)
		}
		BogusSubnets = append(BogusSubnets, n)
	}
}
