package main

// Add predefined in the config blacklist records to the domains list on the redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func InjectBlacklist() {
	log.Info("inject domains from blacklist")

	// create a new redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			CFG.Parse.Redis.Host,
			CFG.Parse.Redis.Port),
		Password: CFG.Parse.Redis.Password,
		DB:       CFG.Parse.Redis.Db,
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Panicln("close redis conn", err)
		}
	}()

	ctx := context.Background()
	rb := NewRBuff("domains")
	defer func() { rb.send(rdb, ctx) }()

	// toss the domains from the config to the redis list
	for _, d := range CFG.BListDomains {
		if domain := sanitizeDomain(d); domainIsOk(domain) {
			rb.add("domains", domain)
		}
	}
}
