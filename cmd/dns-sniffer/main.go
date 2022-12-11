package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/florianl/go-nfqueue"
	"github.com/go-redis/redis/v8"

	log "github.com/sirupsen/logrus"
)

// regQ register queue processing function
func regQ(q uint16, rdb *redis.Client, ac *AllowsCache) {
	config := nfqueue.Config{
		NfQueue:      q,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 50 * time.Millisecond,
	}

	log.Debugf("run hook on %d queue...", q)
	nf, err := nfqueue.Open(&config)
	if err != nil {
		log.Panicln("nfqueue opening", err)
	}

	// add queue closing to teardowns
	td = append(td, func() {
		if e := nf.Close(); e != nil {
			log.Errorln("nfqueue closing", err)
		}
	})

	// create nf_queue processing function
	fn := func(a nfqueue.Attribute) int {
		packetsHook(a, rdb, ac)
		err = nf.SetVerdictWithMark(*a.PacketID, nfqueue.NfRepeat, CFG.MarkDone)
		if err != nil {
			if !strings.Contains(err.Error(), "timeout") {
				log.Panicln("set nfqueue verdict", err)
			}
		}
		return 0
	}

	// register queue processing function
	log.Debugln("register hook func...")
	err = nf.RegisterWithErrorFunc(
		context.Background(), fn, func(e error) int { return 1 })
	if err != nil {
		log.Panicln("nfqueue register fn", err)
	}
}

func main() {
	CFG = initConfig()

	// set logging params
	if CFG.LogLevel != "debug" {
		logLevel, err := log.ParseLevel(CFG.LogLevel)
		if err != nil {
			log.Panic("parsing LogLevel error")
		}
		log.SetLevel(logLevel)
	}
	defer log.Warn("-- done --")

	// setup Redis connection
	log.Debugln("make redis client...")
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

	// create captured answers caching
	allows := NewAllows(rdb, "allows")

	// recheck records live times and cleanup
	go allows.RunSanitizer(time.Second * 5)

	// run nf queue processors
	for _, qn := range CFG.Queues {
		go regQ(qn, rdb, allows)
	}
	<-context.Background().Done()
}
