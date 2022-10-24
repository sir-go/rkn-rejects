package main

import (
	"io"

	log "github.com/sirupsen/logrus"
)

type LogBuff struct {
	records []string
	w       io.Writer
}

func (l *LogBuff) add(msg string) {
	for _, r := range l.records {
		if r == msg {
			return
		}
	}
	l.records = append(l.records, msg)
}

func (l *LogBuff) flush() {
	for _, r := range l.records {
		if _, err := io.WriteString(l.w, r+"\n"); err != nil {
			log.Panicln("can't write log buffer", err)
		}
	}
	l.records = []string{}
}
