package main

import (
	"encoding/xml"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tiaguinho/gosoap"
)

type (
	TimestampMs struct {
		time.Time
	}

	SResVersion struct {
		LastDumpDate         *TimestampMs `xml:"lastDumpDate"`
		LastDumpDateUrgently *TimestampMs `xml:"lastDumpDateUrgently"`
		WebServiceVersion    string       `xml:"webServiceVersion"`
		DumpFormatVersion    string       `xml:"dumpFormatVersion"`
		DocVersion           string       `xml:"docVersion"`
	}
)

func (p *TimestampMs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int64
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	p.Time = time.Unix(v/1000, 0)
	return nil
}

func GetDumpVersion() *SResVersion {
	log.Println("get actual RKN service versions...")
	httpClient := &http.Client{Timeout: CFG.Web.TcpTimeout.Duration}
	defer func() { httpClient.CloseIdleConnections() }()

	soap, err := gosoap.SoapClient(CFG.Web.SoapUrl, httpClient)
	if err != nil {
		log.Panicln("can't make soap client", err)
	}

	soapResp := new(gosoap.Response)
	err = retry(
		CFG.Web.Attempts,
		time.Second*CFG.Web.TcpTimeout.Duration,
		func() error {
			soapResp, err = soap.Call("getLastDumpDateEx", nil)
			return err
		})
	if err != nil {
		log.Panicln("can't call soap method getLastDumpDateEx", err)
	}

	v := new(SResVersion)
	if err = soapResp.Unmarshal(v); err != nil {
		log.Panicln("can't unmarshal soap response", err)
	}
	log.Info(v)
	return v
}

func CheckVersions() {
	log.Info("check versions")
	version := GetDumpVersion()
	if version.DocVersion != CFG.ActualVersions.Doc ||
		version.DumpFormatVersion != CFG.ActualVersions.Dump ||
		version.WebServiceVersion != CFG.ActualVersions.Service {
		log.Errorln("versions mismatch, in config:", CFG.ActualVersions,
			"actual doc:", CFG.Web.DocUrl)
		log.Warn("--break--")
		os.Exit(1)
	}
}
