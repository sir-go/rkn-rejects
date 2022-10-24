package main

import (
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	verdict struct {
		opened bool
		hash   string
		target string
		raw    []byte
	}
)

var reTarget = regexp.MustCompile(`^[^#]((.*)\|)?(.*://)?(.*)`)

func (v *verdict) dump(vDir string) {
	err := ioutil.WriteFile(path.Join(vDir, v.hash), v.raw, 0666)
	if err != nil {
		log.Panicln("dump verdict", err)
	}
}

func check(target string, timeout time.Duration) (v verdict) {
	var (
		err  error
		resp *http.Response
		bb   []byte
	)
	reG := reTarget.FindStringSubmatch(target)
	v = verdict{opened: false}
	if reG == nil {
		return
	}
	hash := reG[2]
	//proto := reG[3]
	href := reG[4]

	client := http.Client{Timeout: timeout}

	for _, proto := range []string{"http", "https"} {
		if resp, err = client.Get(proto + "://" + href); err == nil {
			v = verdict{true, hash, target, []byte{}}
			bb, err = ioutil.ReadAll(resp.Body)
			if len(bb) > 100 {
				bb = bb[:99]
			}
			v.raw = append([]byte(resp.Status+"\n"+href+"\n"), bb...)
			return
		}
	}
	return
}

func Checker(wg *sync.WaitGroup, timeout time.Duration, targets <-chan string,
	verdicts chan<- verdict) {
	for t := range targets {
		verdicts <- check(t, timeout)
		time.Sleep(CFG.Sleeps)
	}
	wg.Done()
}
