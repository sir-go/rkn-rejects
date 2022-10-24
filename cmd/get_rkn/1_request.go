package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/tiaguinho/gosoap"
	"golang.org/x/text/encoding/charmap"

	log "github.com/sirupsen/logrus"
)

type (
	SResReq struct {
		Result  bool   `xml:"result"`
		Comment string `xml:"resultComment"`
		Code    string `xml:"code"`
	}
)

func genRequest() {
	log.Info("gen rewuest")
	var err error
	req := fmt.Sprintf(
		`<?xml version="1.0" encoding="windows-1251"?>
		<request>
		  <requestTime>%s.000+03:00</requestTime>
		  <operatorName>%s</operatorName>
		  <inn>%s</inn>
		  <ogrn>%s</ogrn>
		  <email>%s</email>
		</request>`,
		time.Now().Format("2006-01-02T15:04:05"),
		CFG.Req.Operator.Name,
		CFG.Req.Operator.INN,
		CFG.Req.Operator.OGRN,
		CFG.Req.Operator.Email)
	if req, err = charmap.Windows1251.NewEncoder().String(req); err != nil {
		log.Panicln("encode reg request", err)
	}
	dumpTo(&CFG.Req.File, req, "dump request XML to:")
}

func sign() {
	log.Info("sign request")
	//goland:noinspection SpellCheckingInspection
	cmd := exec.Command(CFG.Sign.Script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debug(cmd)
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}
}

func sendRequest() (taskId string) {
	log.Info("send request")
	httpClient := &http.Client{Timeout: CFG.Web.TcpTimeout.Duration}
	defer func() { httpClient.CloseIdleConnections() }()

	soap, err := gosoap.SoapClient(CFG.Web.SoapUrl, httpClient)
	if err != nil {
		log.Panicln("can't make soap client", err)
	}

	req, err := ioutil.ReadFile(CFG.Req.File)
	if err != nil {
		log.Panicln("can't read req file", CFG.Req.File)
	}

	signed, err := ioutil.ReadFile(CFG.Sign.File)
	if err != nil {
		log.Panicln("can't read signed file", CFG.Sign.File)
	}

	soapResp := new(gosoap.Response)
	err = retry(
		CFG.Web.Attempts,
		time.Second*CFG.Web.TcpTimeout.Duration,
		func() error {
			soapResp, err = soap.Call("sendRequest", gosoap.Params{
				"requestFile":       base64.StdEncoding.EncodeToString(req),
				"signatureFile":     base64.StdEncoding.EncodeToString(signed),
				"dumpFormatVersion": CFG.ActualVersions.Dump,
			})
			return err
		})
	if err != nil {
		log.Panicln("can't call soap method sendRequest", err)
	}

	v := new(SResReq)
	if err = soapResp.Unmarshal(v); err != nil {
		log.Panicln("can't unmarshal soap response", err)
	}
	log.Debug(v)
	if !v.Result {
		log.Error("returned not Ok result")
		log.Warn("--break--")
		os.Exit(1)
	}

	return v.Code
}
