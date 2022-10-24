package main

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type AllowsCache struct {
	rdb *redis.Client
	key string
}

func NewAllows(rdb *redis.Client, key string) *AllowsCache {
	return &AllowsCache{rdb, key}
}

func (ac *AllowsCache) sanitize(score int64) error {
	return ac.rdb.ZRemRangeByScore(context.Background(),
		ac.key, "0", strconv.FormatInt(score, 10)).Err()
}

func (ac *AllowsCache) Add(val string, score int64) error {
	z := redis.Z{Score: float64(score), Member: val}
	return ac.rdb.ZAddArgs(context.Background(),
		ac.key, redis.ZAddArgs{
			Members: []redis.Z{z},
			GT:      true,
		}).Err()
}

func (ac *AllowsCache) Has(val string) bool {
	return ac.rdb.ZScore(context.Background(), ac.key, val).Val() != 0
}

func (ac *AllowsCache) Del(val string) error {
	return ac.rdb.ZRem(context.Background(), ac.key, val).Err()
}

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
