package main

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// A buffer for requests to the Redis server to reduce the load of the parsing processes

type (
	// RBuff contains a mapped requests (list name as a key and a slice of requests as a value)
	// and a total count of stored requests.
	// When a buffer is full it's flushes all of stored requests to the Redis and resets itself
	RBuff struct {
		data  map[string][]interface{}
		count int
	}
)

// NewRBuff creates a buffer mapped with given lists
func NewRBuff(lists ...string) *RBuff {
	rb := &RBuff{data: make(map[string][]interface{})}
	for _, ln := range lists {
		rb.data[ln] = make([]interface{}, 0)
	}
	return rb
}

// add a `val` request to the `list`
func (r *RBuff) add(list string, val interface{}) {
	r.data[list] = append(r.data[list], val)
	r.count++
}

// send flushes all stored requests to the redis via `rdb` connection and recreates
func (r *RBuff) send(rdb *redis.Client, ctx context.Context) {
	keys := make([]string, 0)
	for key, values := range r.data {
		keys = append(keys, key)
		if len(values) > 0 {
			if err := rdb.SAdd(ctx, key, values...).Err(); err != nil {
				log.Panicln("add redis val", err)
			}
		}
	}
	*r = *NewRBuff(keys...)
}
