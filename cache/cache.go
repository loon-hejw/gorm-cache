package cache

import (
	"sync"

	"github.com/loon-hejw/gorm-cache/cache/local"
	"github.com/loon-hejw/gorm-cache/cache/redis"
	"github.com/loon-hejw/gorm-cache/cache/types"
	"github.com/loon-hejw/gorm-cache/config"
)

type Cache interface {
	Set(key string, val interface{}) error
	Get(key string, val interface{}) error
	Del(key string) error
}

type GetCache func(cfg *config.Cache) Cache

var (
	cache    Cache
	once     = sync.Once{}
	cacheMap = make(map[string]GetCache)
)

func Register(key string, f GetCache) {
	cacheMap[key] = f
}

func Init(cfg *config.Cache) error {

	switch cfg.Type {
	case types.CacheRedis:
		cache = redis.New(types.GetRedis(cfg))
	case types.CacheLocal:
		cache = local.New(types.GetLocal(cfg))
	default:
		if f, ok := cacheMap[cfg.Type]; ok {
			cache = f(cfg)
		}
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
