package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type LogFormat struct{}

var CFG *Cfg

func (l *LogFormat) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%-29s [%s] %-15s: %v\n",
		entry.Time.Format("2006/01/02 15:04:05.00.999999"),
		entry.Level.String()[0:1],
		fmt.Sprint(path.Base(entry.Caller.File), ":", entry.Caller.Line),
		entry.Message),
	), nil
}

func initLogging() {
	log.SetReportCaller(true)
	log.SetFormatter(&LogFormat{})
	log.SetLevel(log.DebugLevel)
	if CFG.LogLevel != "debug" {
		logLevel, err := log.ParseLevel(CFG.LogLevel)
		if err != nil {
			log.Panicln("parsing LogLevel error", err)
		}
		log.SetLevel(logLevel)
	}
}

func main() {
	CFG = initConfig()
	initLogging()
	defer log.Println("--done--")

	var dumpData []byte
	if CFG.Parse.FromDump != nil {
		log.Infoln("load dump", *CFG.Parse.FromDump)
		fd, err := os.Open(*CFG.Parse.FromDump)
		if err != nil {
			log.Panicln("open", *CFG.Parse.FromDump, err)
		}
		defer func() {
			if err = fd.Close(); err != nil {
				log.Panicln("can't close", *CFG.Parse.FromDump, err)
			}
		}()

		if dumpData, err = ioutil.ReadAll(fd); err != nil {
			log.Panicln("read dump", err)
		}
	} else {
		CheckVersions()
		genRequest()
		sign()
		taskId := sendRequest()
		dumpData = getResult(taskId)
	}

	dataReader, dataSize := ReadZip(dumpData)
	defer func() {
		if err := dataReader.Close(); err != nil {
			log.Panicln("close zip reader", err)
		}
	}()

	Parse(dataReader, dataSize)
	InjectBlacklist()
	ip2nft()
}
