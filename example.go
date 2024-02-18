package main

import (
	"context"

	"github.com/loon-hejw/gorm-cache/cache"
	"github.com/loon-hejw/gorm-cache/config"
	"gorm.io/gorm"
)

func main() {

	var dial gorm.Dialector
	gdb, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := cache.Init(&config.Cache{}); err != nil {
		panic(err)
	}

	gdb.Use(NewCachePlugin("my-plugin", cache.Session()))

	ctx := context.Background()
	// use cache
	_ = gdb.Session(&gorm.Session{Context: ctx})

	// non use cache
	_ = gdb.Session(&gorm.Session{Context: Unuse(ctx)})
}
