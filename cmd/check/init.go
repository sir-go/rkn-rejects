package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"

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
	CFG   *Cfg
	lBuff LogBuff
)

func InitInterrupt(tearDown func()) {
	log.Infoln("-- start --")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func(c chan os.Signal) {
		for {
			select {
			case <-c:
				tearDown()
				log.Infoln("-- stop --")
				os.Exit(137)
			}
		}
	}(c)
}

func Stop() {
	lBuff.flush()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetReportCaller(true)
	log.SetFormatter(&LogFormat{})
	log.SetLevel(log.DebugLevel)
	InitInterrupt(Stop)
	CFG = initConfig()

	if CFG.LogLevel != "debug" {
		logLevel, err := log.ParseLevel(CFG.LogLevel)
		if err != nil {
			log.Panicln("parsing LogLevel error", err)
		}
		log.SetLevel(logLevel)
	}
}
