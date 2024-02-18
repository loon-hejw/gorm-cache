package config

import (
	"time"
)

const (
	CachePlugin = "gorm:cache"
)

type Cache struct {
	Type  string `yaml:"type,omitempty"`
	Redis Redis  `yaml:"redis,omitempty"`
	Local struct {
		DefaultExpiration time.Duration `yaml:"defaultExpiration,omitempty"`
		CleanupInterval   time.Duration `yaml:"cleanupInterval,omitempty"`
	} `yaml:"local,omitempty"`
	Expiration time.Duration `yaml:"expiration,omitempty"`
}

type Redis struct {
	Addr         string        `yaml:"addr"`
	Username     string        `yaml:"username,omitempty"`
	Password     string        `yaml:"password,omitempty"`
	DB           int           `yaml:"db"`
	DialTimeout  time.Duration `yaml:"dialtimeout,omitempty"`
	ReadTimeout  time.Duration `yaml:"readtimeout,omitempty"`
	WriteTimeout time.Duration `yaml:"writetimeout,omitempty"`
	Pool         struct {
		MaxIdle     int           `yaml:"maxidle,omitempty"`
		MaxActive   int           `yaml:"maxactive,omitempty"`
		IdleTimeout time.Duration `yaml:"idletimeout,omitempty"`
	} `yaml:"pool,omitempty"`
}
