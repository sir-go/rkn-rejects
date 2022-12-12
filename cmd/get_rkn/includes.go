package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// sanitizeConfLine removes from the given string spaces and comments
func sanitizeConfLine(s string) string {
	if s == "" {
		return ""
	}
	di := strings.IndexRune(s, '#')
	if di == 0 {
		return ""
	}
	if di < 0 {
		return strings.TrimSpace(s)
	} else {
		return strings.TrimSpace(s[:di])
	}
}

// includesReadDomains reads a file with domains, recode them and pushes to a string slice
func includesReadDomains(path string) (d []string) {
	fd, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Panicln("can't open domains file", err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			log.Errorln("can't close domains file", err)
		}
	}()
	sc := bufio.NewScanner(fd)
	d = make([]string, 0)
	for sc.Scan() {
		if l := sanitizeConfLine(sc.Text()); l != "" {
			d = append(d, forcePUnicode(l))
		}
	}
	return
}
