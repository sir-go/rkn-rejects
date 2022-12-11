package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var lBuffPath string

func prepareLogBuffer() (logWriter *io.PipeWriter) {
	logWriter = log.StandardLogger().Writer()
	if lBuffPath != "" && lBuffPath != "stdout" {
		fd, err := os.Create(filepath.Clean(lBuffPath))
		if err != nil {
			log.Panicln("can't create log buffer out file", err)
		}
		defer func() {
			if err := fd.Close(); err != nil {
				log.Errorln("can't close log buffer out file", err)
			}
		}()
		lBuff.w = fd
	} else {
		lBuff.w = logWriter
	}
	return
}

func getTargets() []string {
	// init redis db connection
	rdb := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", CFG.Redis.Host, CFG.Redis.Port),
		Password:    CFG.Redis.Password,
		DB:          CFG.Redis.Db,
		DialTimeout: CFG.Redis.TimeoutConn,
		ReadTimeout: CFG.Redis.TimeoutRead,
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Errorln("redis connection closing", err.Error())
		}
	}()

	ctx := context.Background()
	checkRecords, err := rdb.SMembers(ctx, CFG.Redis.SetKey).Result()
	if err != nil {
		log.Panicln("redis get check members:", err)
	}

	if CFG.Max > 0 && len(checkRecords) > CFG.Max {
		checkRecords = checkRecords[:CFG.Max]
	}
	return checkRecords
}

func main() {
	CFG = initConfig()

	if CFG.LogLevel != "debug" {
		logLevel, err := log.ParseLevel(CFG.LogLevel)
		if err != nil {
			log.Panicln("parsing LogLevel error", err)
		}
		log.SetLevel(logLevel)
	}

	var (
		err error
		wg  sync.WaitGroup
	)
	logWriter := prepareLogBuffer()

	// init channels and start workers
	wg.Add(CFG.Workers)

	targets := make(chan string)
	verdicts := make(chan verdict)
	waitCh := make(chan struct{})

	for i := CFG.Workers; i > 0; i-- {
		go Checker(&wg, CFG.Timeout, targets, verdicts)
	}

	checkRecords := getTargets()

	// push records to targets channel
	go func() {
		for _, chR := range checkRecords {
			targets <- chR
		}
		close(targets)

		wg.Wait()
		close(waitCh)
	}()

	progress := 0
	progressOpened := 0
	tick := time.NewTicker(CFG.LogInterval)

	// main loop
	for {
		select {

		// processing verdicts
		case v := <-verdicts:
			progress++
			if !v.opened {
				continue
			}

			// push the verdict to the logging buffer if a target has been opened
			if lBuff.w != logWriter {
				_, err = io.WriteString(lBuff.w, v.target+"\n")
				if err != nil {
					log.Panicln("can't write to log buffer writer", err)
				}
			} else {
				lBuff.add(v.target)
			}
			progressOpened++

			// dump verdict to its file for the farther analyzing
			v.dump(CFG.VerdictsOutDir)

		//	all workers are done
		case <-waitCh:
			log.Infof("%d (! %d)", progress, progressOpened)
			lBuff.flush()
			log.Infof("-- done --")
			return

		//	update the progress
		case <-tick.C:
			log.Infof("%d (! %d)", progress, progressOpened)
		}
	}
}
