package main

import (
	"fmt"
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
	CFG *Cfg     // global config pointer
	td  []func() // teardown callbacks array
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
}
