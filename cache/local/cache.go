package local

import (
	"fmt"

	"reflect"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/hejw123/gorm-cache/cache/types"
)

type inmemory struct {
	cache      *cache.Cache
	expiration time.Duration
}

func New(conf types.Local) *inmemory {
	return &inmemory{
		cache: cache.New(
			conf.DefaultExpiration,
			conf.CleanupInterval,
		),
		expiration: conf.Expiration,
	}
}

func (c *inmemory) Set(key string, val interface{}) error {
	if err := c.cache.Replace(key, val, c.expiration); err != nil {
		return c.cache.Add(key, val, c.expiration)
	}
	return nil
}

func (c *inmemory) Get(key string, val interface{}) error {

	v, ok := c.cache.Get(key)
	if !ok {
		return fmt.Errorf("not found")
	}

	dst := reflect.ValueOf(val).Elem()
	src := reflect.ValueOf(v)

	if dst.Kind() != src.Kind() || dst.Type().String() != src.Type().String() {
		return fmt.Errorf("val kind eq cache val kind")
	}

	dst.Set(src)
	return nil
}

func (c *inmemory) Del(key string) error {
	c.cache.Delete(key)
	return nil
}
