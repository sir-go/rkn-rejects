package main

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type (
	RBuff struct {
		data  map[string][]interface{}
		count int
	}
)

func NewRBuff(lists ...string) *RBuff {
	rb := &RBuff{data: make(map[string][]interface{})}
	for _, ln := range lists {
		rb.data[ln] = make([]interface{}, 0)
	}
	return rb
}

func (r *RBuff) add(list string, val interface{}) {
	r.data[list] = append(r.data[list], val)
	r.count++
}

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
