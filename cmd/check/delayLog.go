package main

// Buffered logs storage. Collects the log records and flushes them to the writer.

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// LogBuff stores logging records and periodically flushes them to the writer
type LogBuff struct {
	records []string
	w       io.Writer
}

// add pushes a log message to the buffer records array
func (l *LogBuff) add(msg string) {
	for _, r := range l.records {
		if r == msg {
			return
		}
	}
	l.records = append(l.records, msg)
}

// flush writes all of the stored records to the writer
func (l *LogBuff) flush() {
	for _, r := range l.records {
		if _, err := io.WriteString(l.w, r+"\n"); err != nil {
			log.Panicln("can't write log buffer", err)
		}
	}
	l.records = []string{}
}
