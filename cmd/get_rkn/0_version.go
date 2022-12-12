package main

// API and a documentation versions checking

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

	// SResVersion - `getLastDumpDateEx` response structure
	SResVersion struct {
		LastDumpDate         *TimestampMs `xml:"lastDumpDate"`
		LastDumpDateUrgently *TimestampMs `xml:"lastDumpDateUrgently"`
		WebServiceVersion    string       `xml:"webServiceVersion"`
		DumpFormatVersion    string       `xml:"dumpFormatVersion"`
		DocVersion           string       `xml:"docVersion"`
	}
)

// UnmarshalXML parses an XML element to a time.Time contained struct
func (p *TimestampMs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int64
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	p.Time = time.Unix(v/1000, 0)
	return nil
}

// GetDumpVersion fetches an XML with API version info and returns structured data
func GetDumpVersion() *SResVersion {
	log.Println("get actual RKN service versions...")
	httpClient := &http.Client{Timeout: CFG.Web.TcpTimeout}
	defer func() { httpClient.CloseIdleConnections() }()

	// create soap client
	soap, err := gosoap.SoapClient(CFG.Web.SoapUrl, httpClient)
	if err != nil {
		log.Panicln("can't make soap client", err)
	}

	// make request & get response
	soapResp := new(gosoap.Response)
	err = retry(
		CFG.Web.Attempts,
		time.Second*CFG.Web.TcpTimeout,
		func() error {
			soapResp, err = soap.Call("getLastDumpDateEx", nil)
			return err
		})
	if err != nil {
		log.Panicln("can't call soap method getLastDumpDateEx", err)
	}

	// decode the answer
	v := new(SResVersion)
	if err = soapResp.Unmarshal(v); err != nil {
		log.Panicln("can't unmarshal soap response", err)
	}
	log.Info(v)
	return v
}

// CheckVersions compares fetched versions and given in the config.
// If versions dont match - throw an error
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
