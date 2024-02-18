package types

import (
	"time"

	"github.com/loon-hejw/gorm-cache/config"
)

const (
	CacheRedis = "redis"
	CacheLocal = "local"
)

var (
	DefaultLocal = Local{
		Config: Config{
			Expiration: 5 * time.Minute,
		},
		DefaultExpiration: 10 * time.Minute,
		CleanupInterval:   10 * time.Minute,
	}
	DefaultRedis = Redis{
		Config: Config{
			Expiration: 5 * time.Minute,
		},
		Dialtimeout:  5 * time.Second,
		Readtimeout:  3 * time.Second,
		Writetimeout: 3 * time.Second,
	}
)

type Config struct {
	Expiration time.Duration
}

type Redis struct {
	Config
	Address         string
	DataBase        int
	User            string
	Pass            string
	Dialtimeout     time.Duration
	Readtimeout     time.Duration
	Writetimeout    time.Duration
	PoolMaxidle     int
	PoolMaxactive   int
	PoolIdletimeout time.Duration
}

type Local struct {
	Config
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func GetRedis(cfg *config.Cache) Redis {

	if cfg == nil {
		return DefaultRedis
	}

	defaultRedis := DefaultRedis
	if cfg.Expiration != 0 {
		defaultRedis.Expiration = cfg.Expiration
	}

	if cfg.Redis.DialTimeout != 0 {
		defaultRedis.Dialtimeout = cfg.Redis.DialTimeout
	}

	if cfg.Redis.ReadTimeout != 0 {
		defaultRedis.Readtimeout = cfg.Redis.ReadTimeout
	}

	if cfg.Redis.WriteTimeout != 0 {
		defaultRedis.Writetimeout = cfg.Redis.WriteTimeout
	}

	if cfg.Redis.Pool.IdleTimeout != 0 {
		defaultRedis.PoolIdletimeout = cfg.Redis.Pool.IdleTimeout
	}

	if cfg.Redis.Pool.MaxActive != 0 {
		defaultRedis.PoolMaxactive = cfg.Redis.Pool.MaxActive
	}

	if cfg.Redis.Pool.MaxIdle != 0 {
		defaultRedis.PoolMaxidle = cfg.Redis.Pool.MaxIdle
	}

	defaultRedis.Address = cfg.Redis.Addr
	defaultRedis.DataBase = cfg.Redis.DB
	defaultRedis.User = cfg.Redis.Username
	defaultRedis.Pass = cfg.Redis.Password
	return defaultRedis
}

func GetLocal(cfg *config.Cache) Local {

	if cfg == nil {
		return DefaultLocal
	}

	defaultLocal := DefaultLocal
	if cfg.Expiration != 0 {
		defaultLocal.Expiration = cfg.Expiration
	}

	if cfg.Local.DefaultExpiration != 0 {
		defaultLocal.DefaultExpiration = cfg.Local.DefaultExpiration
	}

	if cfg.Local.CleanupInterval != 0 {
		defaultLocal.CleanupInterval = cfg.Local.CleanupInterval
	}
	return defaultLocal
}
