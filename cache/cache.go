package cache

import (
	"sync"

	"github.com/spf13/viper"

	"github.com/hejw123/gorm-cache/cache/local"
	"github.com/hejw123/gorm-cache/cache/redis"
	"github.com/hejw123/gorm-cache/cache/types"
	"github.com/hejw123/gorm-cache/config"
)

type Cache interface {
	Set(key string, val interface{}) error
	Get(key string, val interface{}) error
	Del(key string) error
}

var (
	cache Cache
	once  = sync.Once{}
)

func Init() error {

	cfg := &config.Cache{}
	if err := viper.UnmarshalKey("cache", &cfg); err != nil {
		return err
	}
	switch cfg.Type {
	case types.CacheRedis:
		cache = redis.New(types.GetRedis(cfg))
	case types.CacheLocal:
		cache = local.New(types.GetLocal(cfg))
	}
	return nil
}

func Session() Cache {
	once.Do(func() {
		if cache != nil {
			return
		}
		cache = local.New(types.DefaultLocal)
	})
	return cache
}
