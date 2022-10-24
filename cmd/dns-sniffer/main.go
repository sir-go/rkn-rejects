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

	td = append(td, func() {
		if e := nf.Close(); e != nil {
			log.Errorln("nfqueue closing", err)
		}
	})

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
	log.Debugln("register hook func...")
	err = nf.RegisterWithErrorFunc(
		context.Background(), fn, func(e error) int { return 1 })
	if err != nil {
		log.Panicln("nfqueue register fn", err)
	}
}

func main() {
	defer log.Warn("-- done --")

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

	allows := NewAllows(rdb, "allows")
	go allows.RunSanitizer(time.Second * 5)

	for _, qn := range CFG.Queues {
		go regQ(qn, rdb, allows)
	}
	<-context.Background().Done()
}
