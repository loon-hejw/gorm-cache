gorm-cache aims to provide a look-aside, almost-no-code-modification cache solution for gorm users.

gorm-cache 旨在为gorm用户提供一个即插即用的旁路缓存解决方案。

```go
package main

import (
	"context"

	"github.com/hejw123/gorm-cache/cache"
	"gorm.io/gorm"
)

func main() {

	var dial gorm.Dialector

	// plugin = NewCachePlugin(cfg.Prefix, cache.Session())
	gdb, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	gdb.Use(NewCachePlugin("my-plugin", cache.Session()))

	ctx := context.Background()

	// use cache
	_ = gdb.Session(&gorm.Session{Context: ctx})

	// non use cache
	_ = gdb.Session(&gorm.Session{Context: NonuseCache(ctx)})
}
```