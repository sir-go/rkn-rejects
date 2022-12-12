package main

// Fetching response as a zip-file bytes

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tiaguinho/gosoap"
)

type (
	// SResResult - `getResult` response structure
	SResResult struct {
		Result            bool    `xml:"result"`
		Code              int     `xml:"resultCode"`
		Comment           *string `xml:"resultComment"`
		ZipArchive        *string `xml:"registerZipArchive"`
		DumpFormatVersion *string `xml:"dumpFormatVersion"`
		OperatorName      *string `xml:"operatorName"`
		INN               *string `xml:"inn"`
	}
)

// getResult continuously tries to fetch the result of the given task ID
// stores the result to a file (if path is presented in the config) and returns a zip dump bytes
func getResult(code string) (zipDump []byte) {
	log.Info("get result")
	httpClient := &http.Client{Timeout: CFG.Res.GetTimeout}
	defer func() { httpClient.CloseIdleConnections() }()

	soap, err := gosoap.SoapClient(CFG.Web.SoapUrl, httpClient)
	if err != nil {
		log.Panicln("can't make soap client", err)
	}

	v := new(SResResult)
	soapResp := new(gosoap.Response)
	err = retry(
		CFG.Res.Attempts,
		CFG.Res.RetryTimeout,
		func() error {
			soapResp, err = soap.Call("getResult",
				gosoap.Params{"code": code})
			if err != nil {
				return err
			}
			if err = soapResp.Unmarshal(v); err != nil {
				return err
			}
			if v.Result {
				return nil
			}

			if v.Comment == nil {
				return errors.New("no comment")
			}

			if v.Code == 0 || v.Code == -10 {
				return errors.New(*v.Comment)
			}
			log.Error(*v.Comment)
			log.Warn("--break--")
			os.Exit(1)
			return nil
		})
	if err != nil {
		log.Panicln("can't call soap method getResult", err)
	}

	zipDump, err = base64.StdEncoding.DecodeString(*v.ZipArchive)
	if err != nil {
		log.Panicln("decode dump", err)
	}

	dumpTo(CFG.Res.DumpTo, zipDump, "dump saved to:")
	return
}
