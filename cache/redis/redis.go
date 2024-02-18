package redis

import (
	"context"
	"time"

	"github.com/loon-hejw/gorm-cache/cache/types"

	rv9 "github.com/redis/go-redis/v9"
)

var ctx = context.TODO()

type redis struct {
	client     *rv9.Client
	expiration time.Duration
}

func New(conf types.Redis) *redis {
	rdb := rv9.NewClient(&rv9.Options{
		Addr:           conf.Address,
		Username:       conf.User,
		Password:       conf.Pass,
		DB:             conf.DataBase,
		PoolTimeout:    conf.PoolIdletimeout,
		MaxIdleConns:   conf.PoolMaxidle,
		MaxActiveConns: conf.PoolMaxactive,
		DialTimeout:    conf.Dialtimeout,
		ReadTimeout:    conf.Readtimeout,
		WriteTimeout:   conf.Writetimeout,
	})
	return &redis{
		client:     rdb,
		expiration: conf.Expiration,
	}
}

func (r *redis) Set(key string, val interface{}) error {
	return r.client.Set(ctx, key, val, r.expiration).Err()
}

func (r *redis) Get(key string, val interface{}) error {
	return r.client.Get(ctx, key).Scan(val)
}

func (r *redis) Del(key string) error {
	return r.client.Del(ctx, key).Err()
}
