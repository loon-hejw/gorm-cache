package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"sync"

	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"

	store "github.com/hejw123/gorm-cache/cache"

	json "github.com/json-iterator/go"
)

type NonUsekey bool

var noCache NonUsekey

const (
	CachePlugin = "gorm:cache"
)

type cachePlugin struct {
	prefix string
	store  store.Cache
	tables map[string]map[string]struct{}
	lock   sync.Mutex
}

func NewCachePlugin(prefix string, store store.Cache) gorm.Plugin {
	return &cachePlugin{
		prefix: prefix,
		store:  store,
		tables: make(map[string]map[string]struct{}),
		lock:   sync.Mutex{},
	}
}

func (c *cachePlugin) Name() string {
	return CachePlugin
}

func (c *cachePlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Delete().After("*").Register("cache:clear", c.ClearCache); err != nil {
		return err
	}
	if err := db.Callback().Update().After("*").Register("cache:clear", c.ClearCache); err != nil {
		return err
	}
	if err := db.Callback().Create().After("*").Register("cache:clear", c.ClearCache); err != nil {
		return err
	}
	return db.Callback().Query().Replace("gorm:query", c.Query)
}

func (c *cachePlugin) ClearCache(tx *gorm.DB) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.tables, tx.Statement.Table)
}

func (c *cachePlugin) Query(tx *gorm.DB) {

	// 调用 Gorm的方法生产SQL
	callbacks.BuildQuerySQL(tx)
	sql := tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
	key := fmt.Sprintf("%s%s", c.prefix, generateKey(sql))
	ctx := tx.Statement.Context

	// 查询缓存数据
	if c.exist(ctx, tx.Statement.Table, key) {
		if err := c.QueryCache(ctx, key, tx); err == nil {
			return
		}
	}

	// 查询数据库
	c.QueryDB(tx)
	if tx.Error != nil {
		return
	}

	value, err := json.MarshalToString(tx.Statement.Dest)
	if err != nil {
		return
	}

	// 写入缓存
	if err := c.store.Set(key, value); err != nil {
		tx.Logger.Error(ctx, err.Error())
		return
	}
}

func (c *cachePlugin) exist(ctx context.Context, key, sql string) bool {

	if use, ok := ctx.Value(noCache).(bool); ok && use {
		return false
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.tables[key]; !ok {
		c.tables[key] = make(map[string]struct{})
	}
	_, ok := c.tables[key][sql]
	if !ok {
		c.tables[key][sql] = struct{}{}
	}
	return ok
}

func (c *cachePlugin) QueryCache(ctx context.Context, key string, tx *gorm.DB) error {

	var value string
	if err := c.store.Get(key, &value); err != nil {
		return err
	}

	decode := json.NewDecoder(strings.NewReader(value))
	switch tx.Statement.Dest.(type) {
	case *int64:
		tx.RowsAffected = 1
		decode.UseNumber()
	}
	return decode.Decode(&tx.Statement.Dest)
}

func (p *cachePlugin) QueryDB(tx *gorm.DB) {
	if tx.Error != nil || tx.DryRun {
		return
	}

	rows, err := tx.Statement.ConnPool.QueryContext(tx.Statement.Context, tx.Statement.SQL.String(), tx.Statement.Vars...)
	if err != nil {
		_ = tx.AddError(err)
		return
	}

	defer func() {
		_ = tx.AddError(rows.Close())
	}()

	gorm.Scan(rows, tx, 0)
}

func generateKey(key string) string {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(key))
	return strconv.FormatUint(hash.Sum64(), 36)
}

func NonuseCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, noCache, true)
}
