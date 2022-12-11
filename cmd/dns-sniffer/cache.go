package main

// Redis-based DNS answers cache,
// stores answers (resolved IP strings arrays) in a ZLists where score is an record expiration time

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// AllowsCache stores a db connection pointer and a db key
type AllowsCache struct {
	rdb *redis.Client
	key string
}

// NewAllows creates a new cache keeper structure
func NewAllows(rdb *redis.Client, key string) *AllowsCache {
	return &AllowsCache{rdb, key}
}

// sanitize removes all zero-score records
func (ac *AllowsCache) sanitize(score int64) error {
	return ac.rdb.ZRemRangeByScore(context.Background(),
		ac.key, "0", strconv.FormatInt(score, 10)).Err()
}

// Add stores a record (IP addr strings with score - expiration time)
func (ac *AllowsCache) Add(val string, score int64) error {
	z := redis.Z{Score: float64(score), Member: val}
	return ac.rdb.ZAddArgs(context.Background(),
		ac.key, redis.ZAddArgs{
			Members: []redis.Z{z},
			GT:      true,
		}).Err()
}

// Has checks if the cache stores the given IP string
func (ac *AllowsCache) Has(val string) bool {
	return ac.rdb.ZScore(context.Background(), ac.key, val).Val() != 0
}

// Del removes the IP address from the cache
func (ac *AllowsCache) Del(val string) error {
	return ac.rdb.ZRem(context.Background(), ac.key, val).Err()
}

// RunSanitizer gets time ticks and cleanup cache (remove all records with expired time)
func (ac *AllowsCache) RunSanitizer(timeout time.Duration) {
	tick := time.NewTicker(timeout + time.Second)
	for {
		select {
		case t := <-tick.C:
			if err := ac.sanitize(t.Unix()); err != nil {
				log.Errorln("redis sanitize: ", err)
			}
		}
	}
}
