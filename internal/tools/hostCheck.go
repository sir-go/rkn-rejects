package tools

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// IsHostDenied checks the hostname and all of the upper domains if any of them is in the denied list
func IsHostDenied(h string, rdb *redis.Client, hostsKey string) bool {
	ctx := context.Background()
	for _, ud := range GetUpperDomains(h) {
		rRes := rdb.SIsMember(ctx, hostsKey, ud)
		if err := rRes.Err(); err != nil {
			log.Errorln("redis sismember", hostsKey, err)
			return false
		}
		if rRes.Val() {
			return true
		}
	}
	return false
}
