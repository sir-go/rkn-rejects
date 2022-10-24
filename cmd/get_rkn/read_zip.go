package main

import (
	"archive/zip"
	"bytes"
	"io"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ReadZip(zipDump []byte) (r io.ReadCloser, size uint64) {
	zipReader, err := zip.NewReader(
		bytes.NewReader(zipDump), int64(len(zipDump)))
	if err != nil {
		log.Panicln("can't read zip dump", err)
	}

	for _, zipFile := range zipReader.File {
		if strings.ToLower(path.Ext(zipFile.Name)) == ".xml" {
			if r, err = zipFile.Open(); err != nil {
				log.Panicln("can't open zip file", err)
			}
			size = zipFile.UncompressedSize64
			return
		}
	}
	return
}
