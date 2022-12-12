package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// dumpTo saves the given data to a file if the path is presented
func dumpTo(path *string, data interface{}, decr string) {
	if path == nil || *path == "" {
		return
	}
	fd, err := os.Create(*path)
	if err != nil {
		log.Panicln("create dump file", err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			log.Errorln("can't close dump file", err)
		}
	}()

	log.Infoln(decr, *path)
	switch data.(type) {
	case string:
		_, err = fd.WriteString(data.(string))
	case []byte:
		_, err = fd.Write(data.([]byte))
	}
	if err != nil {
		log.Panicln("can't write dump", err)
	}
}
